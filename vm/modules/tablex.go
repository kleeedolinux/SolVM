package modules

import (
	"fmt"
	"math/rand"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type TableX struct {
	L *lua.LState
}

func RegisterTableXModule(L *lua.LState) {
	tablex := &TableX{L: L}
	mod := L.NewTable()

	L.SetField(mod, "copy", L.NewFunction(tablex.copy))
	L.SetField(mod, "deepcopy", L.NewFunction(tablex.deepcopy))
	L.SetField(mod, "compare", L.NewFunction(tablex.compare))
	L.SetField(mod, "map", L.NewFunction(tablex.mapTable))
	L.SetField(mod, "filter", L.NewFunction(tablex.filter))
	L.SetField(mod, "reduce", L.NewFunction(tablex.reduce))
	L.SetField(mod, "flatten", L.NewFunction(tablex.flatten))
	L.SetField(mod, "keys", L.NewFunction(tablex.keys))
	L.SetField(mod, "values", L.NewFunction(tablex.values))
	L.SetField(mod, "pretty", L.NewFunction(tablex.pretty))
	L.SetField(mod, "load", L.NewFunction(tablex.load))
	L.SetField(mod, "loadfile", L.NewFunction(tablex.loadfile))

	L.SetField(mod, "slice", L.NewFunction(tablex.slice))
	L.SetField(mod, "concat", L.NewFunction(tablex.concat))
	L.SetField(mod, "partition", L.NewFunction(tablex.partition))
	L.SetField(mod, "rotate", L.NewFunction(tablex.rotate))
	L.SetField(mod, "shuffle", L.NewFunction(tablex.shuffle))

	L.SetField(mod, "map_new", L.NewFunction(tablex.mapNew))
	L.SetField(mod, "set_new", L.NewFunction(tablex.setNew))
	L.SetField(mod, "ordered_map_new", L.NewFunction(tablex.orderedMapNew))

	L.SetField(mod, "array2d_new", L.NewFunction(tablex.array2DNew))
	L.SetField(mod, "array2d_get", L.NewFunction(tablex.array2DGet))
	L.SetField(mod, "array2d_set", L.NewFunction(tablex.array2DSet))
	L.SetField(mod, "array2d_map", L.NewFunction(tablex.array2DMap))
	L.SetField(mod, "array2d_filter", L.NewFunction(tablex.array2DFilter))
	L.SetField(mod, "array2d_transpose", L.NewFunction(tablex.array2DTranspose))

	L.SetField(mod, "permute", L.NewFunction(tablex.permute))
	L.SetField(mod, "combinations", L.NewFunction(tablex.combinations))

	L.SetGlobal("tablex", mod)
}

func (t *TableX) copy(L *lua.LState) int {
	tbl := L.CheckTable(1)
	result := L.NewTable()
	tbl.ForEach(func(key, value lua.LValue) {
		result.RawSet(key, value)
	})
	L.Push(result)
	return 1
}

func (t *TableX) deepcopy(L *lua.LState) int {
	tbl := L.CheckTable(1)
	result := L.NewTable()
	t.deepcopyTable(tbl, result)
	L.Push(result)
	return 1
}

func (t *TableX) deepcopyTable(src, dst *lua.LTable) {
	src.ForEach(func(key, value lua.LValue) {
		switch value.Type() {
		case lua.LTTable:
			newTbl := t.L.NewTable()
			t.deepcopyTable(value.(*lua.LTable), newTbl)
			dst.RawSet(key, newTbl)
		default:
			dst.RawSet(key, value)
		}
	})
}

func (t *TableX) compare(L *lua.LState) int {
	tbl1 := L.CheckTable(1)
	tbl2 := L.CheckTable(2)
	L.Push(lua.LBool(t.compareTables(tbl1, tbl2)))
	return 1
}

func (t *TableX) compareTables(tbl1, tbl2 *lua.LTable) bool {
	if tbl1.Len() != tbl2.Len() {
		return false
	}

	equal := true
	tbl1.ForEach(func(key, value lua.LValue) {
		if !equal {
			return
		}
		value2 := tbl2.RawGet(key)
		if value.Type() != value2.Type() {
			equal = false
			return
		}
		if value.Type() == lua.LTTable {
			equal = t.compareTables(value.(*lua.LTable), value2.(*lua.LTable))
		} else {
			equal = value == value2
		}
	})
	return equal
}

func (t *TableX) mapTable(L *lua.LState) int {
	tbl := L.CheckTable(1)
	fn := L.CheckFunction(2)
	result := L.NewTable()

	tbl.ForEach(func(key, value lua.LValue) {
		L.Push(fn)
		L.Push(value)
		if err := L.PCall(1, 1, nil); err != nil {
			L.RaiseError("error in map function: %v", err)
		}
		result.RawSet(key, L.Get(-1))
		L.Pop(1)
	})

	L.Push(result)
	return 1
}

func (t *TableX) filter(L *lua.LState) int {
	tbl := L.CheckTable(1)
	fn := L.CheckFunction(2)
	result := L.NewTable()

	tbl.ForEach(func(key, value lua.LValue) {
		L.Push(fn)
		L.Push(value)
		if err := L.PCall(1, 1, nil); err != nil {
			L.RaiseError("error in filter function: %v", err)
		}
		if L.ToBool(-1) {
			result.RawSet(key, value)
		}
		L.Pop(1)
	})

	L.Push(result)
	return 1
}

func (t *TableX) reduce(L *lua.LState) int {
	tbl := L.CheckTable(1)
	fn := L.CheckFunction(2)
	var initial lua.LValue
	if L.GetTop() > 2 {
		initial = L.Get(3)
	} else {
		initial = lua.LNil
	}

	acc := initial
	tbl.ForEach(func(_, value lua.LValue) {
		L.Push(fn)
		L.Push(acc)
		L.Push(value)
		if err := L.PCall(2, 1, nil); err != nil {
			L.RaiseError("error in reduce function: %v", err)
		}
		acc = L.Get(-1)
		L.Pop(1)
	})

	L.Push(acc)
	return 1
}

func (t *TableX) flatten(L *lua.LState) int {
	tbl := L.CheckTable(1)
	depth := L.OptInt(2, -1)
	result := L.NewTable()
	t.flattenTable(tbl, result, depth, 0)
	L.Push(result)
	return 1
}

func (t *TableX) flattenTable(src, dst *lua.LTable, maxDepth, currentDepth int) {
	src.ForEach(func(_, value lua.LValue) {
		if value.Type() == lua.LTTable && (maxDepth == -1 || currentDepth < maxDepth) {
			t.flattenTable(value.(*lua.LTable), dst, maxDepth, currentDepth+1)
		} else {
			dst.Append(value)
		}
	})
}

func (t *TableX) keys(L *lua.LState) int {
	tbl := L.CheckTable(1)
	result := L.NewTable()
	tbl.ForEach(func(key, _ lua.LValue) {
		result.Append(key)
	})
	L.Push(result)
	return 1
}

func (t *TableX) values(L *lua.LState) int {
	tbl := L.CheckTable(1)
	result := L.NewTable()
	tbl.ForEach(func(_, value lua.LValue) {
		result.Append(value)
	})
	L.Push(result)
	return 1
}

func (t *TableX) pretty(L *lua.LState) int {
	tbl := L.CheckTable(1)
	indent := L.OptInt(2, 0)
	L.Push(lua.LString(t.prettyTable(tbl, indent)))
	return 1
}

func (t *TableX) prettyTable(tbl *lua.LTable, indent int) string {
	if tbl.Len() == 0 {
		return "{}"
	}

	var result strings.Builder
	indentStr := strings.Repeat("  ", indent)
	result.WriteString("{\n")

	tbl.ForEach(func(key, value lua.LValue) {
		result.WriteString(indentStr)
		result.WriteString("  ")
		result.WriteString(t.prettyValue(key))
		result.WriteString(" = ")
		result.WriteString(t.prettyValue(value))
		result.WriteString(",\n")
	})

	result.WriteString(indentStr)
	result.WriteString("}")
	return result.String()
}

func (t *TableX) prettyValue(value lua.LValue) string {
	switch value.Type() {
	case lua.LTTable:
		return t.prettyTable(value.(*lua.LTable), 0)
	case lua.LTString:
		return fmt.Sprintf("%q", value.String())
	default:
		return value.String()
	}
}

func (t *TableX) load(L *lua.LState) int {
	str := L.CheckString(1)
	fn, err := t.L.LoadString("return " + str)
	if err != nil {
		L.RaiseError("error loading table: %v", err)
	}
	t.L.Push(fn)
	if err := t.L.PCall(0, 1, nil); err != nil {
		L.RaiseError("error executing table: %v", err)
	}
	return 1
}

func (t *TableX) loadfile(L *lua.LState) int {
	filename := L.CheckString(1)
	fn, err := t.L.LoadFile(filename)
	if err != nil {
		L.RaiseError("error loading file: %v", err)
	}
	t.L.Push(fn)
	if err := t.L.PCall(0, 1, nil); err != nil {
		L.RaiseError("error executing file: %v", err)
	}
	return 1
}

func (t *TableX) slice(L *lua.LState) int {
	tbl := L.CheckTable(1)
	start := L.OptInt(2, 1)
	end := L.OptInt(3, tbl.Len())
	step := L.OptInt(4, 1)

	if start < 1 {
		start = 1
	}
	if end > tbl.Len() {
		end = tbl.Len()
	}

	result := L.NewTable()
	for i := start; i <= end; i += step {
		result.Append(tbl.RawGetInt(i))
	}

	L.Push(result)
	return 1
}

func (t *TableX) concat(L *lua.LState) int {
	tbl := L.CheckTable(1)
	sep := L.OptString(2, "")
	result := L.NewTable()

	tbl.ForEach(func(_, value lua.LValue) {
		if value.Type() == lua.LTTable {
			subResult := L.NewTable()
			t.concatTable(value.(*lua.LTable), sep, subResult)
			subResult.ForEach(func(_, v lua.LValue) {
				result.Append(v)
			})
		} else {
			result.Append(value)
		}
	})

	L.Push(result)
	return 1
}

func (t *TableX) concatTable(src *lua.LTable, sep string, dst *lua.LTable) {
	src.ForEach(func(_, value lua.LValue) {
		if value.Type() == lua.LTTable {
			t.concatTable(value.(*lua.LTable), sep, dst)
		} else {
			dst.Append(value)
		}
	})
}

func (t *TableX) partition(L *lua.LState) int {
	tbl := L.CheckTable(1)
	size := L.CheckInt(2)
	result := L.NewTable()

	for i := 1; i <= tbl.Len(); i += size {
		part := L.NewTable()
		for j := 0; j < size && i+j <= tbl.Len(); j++ {
			part.Append(tbl.RawGetInt(i + j))
		}
		result.Append(part)
	}

	L.Push(result)
	return 1
}

func (t *TableX) rotate(L *lua.LState) int {
	tbl := L.CheckTable(1)
	n := L.OptInt(2, 1)
	result := L.NewTable()

	len := tbl.Len()
	if len == 0 {
		L.Push(result)
		return 1
	}

	n = ((n % len) + len) % len
	for i := 1; i <= len; i++ {
		idx := ((i + n - 1) % len) + 1
		result.Append(tbl.RawGetInt(idx))
	}

	L.Push(result)
	return 1
}

func (t *TableX) shuffle(L *lua.LState) int {
	tbl := L.CheckTable(1)
	result := L.NewTable()
	values := make([]lua.LValue, 0, tbl.Len())

	tbl.ForEach(func(_, value lua.LValue) {
		values = append(values, value)
	})

	for i := len(values) - 1; i > 0; i-- {
		j := int(float64(i+1) * rand.Float64())
		values[i], values[j] = values[j], values[i]
	}

	for _, value := range values {
		result.Append(value)
	}

	L.Push(result)
	return 1
}

func (t *TableX) mapNew(L *lua.LState) int {
	result := L.NewTable()
	L.SetMetatable(result, L.GetTypeMetatable("Map"))
	L.Push(result)
	return 1
}

func (t *TableX) setNew(L *lua.LState) int {
	result := L.NewTable()
	L.SetMetatable(result, L.GetTypeMetatable("Set"))
	L.Push(result)
	return 1
}

func (t *TableX) orderedMapNew(L *lua.LState) int {
	result := L.NewTable()
	L.SetMetatable(result, L.GetTypeMetatable("OrderedMap"))
	L.Push(result)
	return 1
}

func (t *TableX) array2DNew(L *lua.LState) int {
	rows := L.CheckInt(1)
	cols := L.CheckInt(2)
	result := L.NewTable()

	mt := L.NewTable()
	L.SetField(mt, "cols", lua.LNumber(cols))
	L.SetMetatable(result, mt)

	for i := 1; i <= rows; i++ {
		row := L.NewTable()

		for j := 1; j <= cols; j++ {
			row.RawSetInt(j, lua.LNil)
		}
		result.RawSetInt(i, row)
	}

	L.Push(result)
	return 1
}

func (t *TableX) array2DGet(L *lua.LState) int {
	arr := L.CheckTable(1)
	row := L.CheckInt(2)
	col := L.CheckInt(3)

	if row < 1 || row > arr.Len() {
		L.RaiseError("row index out of bounds: %d (max: %d)", row, arr.Len())
	}

	mt := L.GetMetatable(arr)
	if mt == lua.LNil {
		L.RaiseError("invalid 2D array: missing metatable")
	}
	cols := int(L.GetField(mt, "cols").(lua.LNumber))

	if col < 1 || col > cols {
		L.RaiseError("column index out of bounds: %d (max: %d)", col, cols)
	}

	rowTbl := arr.RawGetInt(row).(*lua.LTable)
	L.Push(rowTbl.RawGetInt(col))
	return 1
}

func (t *TableX) array2DSet(L *lua.LState) int {
	arr := L.CheckTable(1)
	row := L.CheckInt(2)
	col := L.CheckInt(3)
	value := L.CheckAny(4)

	if row < 1 || row > arr.Len() {
		L.RaiseError("row index out of bounds: %d (max: %d)", row, arr.Len())
	}

	mt := L.GetMetatable(arr)
	if mt == lua.LNil {
		L.RaiseError("invalid 2D array: missing metatable")
	}
	cols := int(L.GetField(mt, "cols").(lua.LNumber))

	if col < 1 || col > cols {
		L.RaiseError("column index out of bounds: %d (max: %d)", col, cols)
	}

	rowTbl := arr.RawGetInt(row).(*lua.LTable)
	rowTbl.RawSetInt(col, value)
	return 0
}

func (t *TableX) array2DMap(L *lua.LState) int {
	arr := L.CheckTable(1)
	fn := L.CheckFunction(2)
	result := L.NewTable()

	for i := 1; i <= arr.Len(); i++ {
		row := arr.RawGetInt(i).(*lua.LTable)
		newRow := L.NewTable()
		for j := 1; j <= row.Len(); j++ {
			L.Push(fn)
			L.Push(row.RawGetInt(j))
			if err := L.PCall(1, 1, nil); err != nil {
				L.RaiseError("error in map function: %v", err)
			}
			newRow.RawSetInt(j, L.Get(-1))
			L.Pop(1)
		}
		result.RawSetInt(i, newRow)
	}

	L.Push(result)
	return 1
}

func (t *TableX) array2DFilter(L *lua.LState) int {
	arr := L.CheckTable(1)
	fn := L.CheckFunction(2)
	result := L.NewTable()

	for i := 1; i <= arr.Len(); i++ {
		row := arr.RawGetInt(i).(*lua.LTable)
		newRow := L.NewTable()
		for j := 1; j <= row.Len(); j++ {
			L.Push(fn)
			L.Push(row.RawGetInt(j))
			if err := L.PCall(1, 1, nil); err != nil {
				L.RaiseError("error in filter function: %v", err)
			}
			if L.ToBool(-1) {
				newRow.Append(row.RawGetInt(j))
			}
			L.Pop(1)
		}
		if newRow.Len() > 0 {
			result.Append(newRow)
		}
	}

	L.Push(result)
	return 1
}

func (t *TableX) array2DTranspose(L *lua.LState) int {
	arr := L.CheckTable(1)
	if arr.Len() == 0 {
		L.Push(arr)
		return 1
	}

	mt := L.GetMetatable(arr)
	if mt == lua.LNil {
		L.RaiseError("invalid 2D array: missing metatable")
	}
	cols := int(L.GetField(mt, "cols").(lua.LNumber))

	result := L.NewTable()

	newMt := L.NewTable()
	L.SetField(newMt, "cols", lua.LNumber(arr.Len()))
	L.SetMetatable(result, newMt)

	for i := 1; i <= cols; i++ {
		row := L.NewTable()
		for j := 1; j <= arr.Len(); j++ {
			row.Append(arr.RawGetInt(j).(*lua.LTable).RawGetInt(i))
		}
		result.RawSetInt(i, row)
	}

	L.Push(result)
	return 1
}

func (t *TableX) permute(L *lua.LState) int {
	tbl := L.CheckTable(1)
	result := L.NewTable()
	values := make([]lua.LValue, 0, tbl.Len())

	tbl.ForEach(func(_, value lua.LValue) {
		values = append(values, value)
	})

	t.permuteValues(values, 0, result)
	L.Push(result)
	return 1
}

func (t *TableX) permuteValues(values []lua.LValue, start int, result *lua.LTable) {
	if start == len(values)-1 {
		perm := t.L.NewTable()
		for _, v := range values {
			perm.Append(v)
		}
		result.Append(perm)
		return
	}

	for i := start; i < len(values); i++ {
		values[start], values[i] = values[i], values[start]
		t.permuteValues(values, start+1, result)
		values[start], values[i] = values[i], values[start]
	}
}

func (t *TableX) combinations(L *lua.LState) int {
	tbl := L.CheckTable(1)
	r := L.CheckInt(2)
	result := L.NewTable()
	values := make([]lua.LValue, 0, tbl.Len())

	tbl.ForEach(func(_, value lua.LValue) {
		values = append(values, value)
	})

	t.combineValues(values, r, 0, make([]lua.LValue, r), result)
	L.Push(result)
	return 1
}

func (t *TableX) combineValues(values []lua.LValue, r, start int, current []lua.LValue, result *lua.LTable) {
	if r == 0 {
		comb := t.L.NewTable()
		for _, v := range current {
			comb.Append(v)
		}
		result.Append(comb)
		return
	}

	for i := start; i <= len(values)-r; i++ {
		current[len(current)-r] = values[i]
		t.combineValues(values, r-1, i+1, current, result)
	}
}
