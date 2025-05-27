-- Table operations
local t1 = {a = 1, b = 2, c = {d = 3, e = 4}}
local t2 = tablex.copy(t1)
print("Copy:", tablex.pretty(t2))

local t3 = tablex.deepcopy(t1)
t3.c.d = 5
print("Deep copy:", tablex.pretty(t3))
print("Original:", tablex.pretty(t1))

-- List operations
local list = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
print("Slice:", tablex.pretty(tablex.slice(list, 2, 5)))
print("Partition:", tablex.pretty(tablex.partition(list, 3)))
print("Rotate:", tablex.pretty(tablex.rotate(list, 3)))
print("Shuffle:", tablex.pretty(tablex.shuffle(list)))

-- Map operations
local mapped = tablex.map(list, function(x) return x * 2 end)
print("Map:", tablex.pretty(mapped))

local filtered = tablex.filter(list, function(x) return x % 2 == 0 end)
print("Filter:", tablex.pretty(filtered))

local sum = tablex.reduce(list, function(acc, x) return acc + x end, 0)
print("Reduce (sum):", sum)

-- 2D Array operations
print("\n2D Array Operations:")
local rows, cols = 3, 3
local arr2d = tablex.array2d_new(rows, cols)

-- Fill the array with values
for i = 1, rows do
    for j = 1, cols do
        local value = (i-1) * cols + j
        tablex.array2d_set(arr2d, i, j, value)
    end
end

-- Print the array
print("2D Array:")
for i = 1, rows do
    local row = {}
    for j = 1, cols do
        row[j] = tablex.array2d_get(arr2d, i, j)
    end
    print(tablex.pretty(row))
end

-- Transpose and print
local transposed = tablex.array2d_transpose(arr2d)
print("\nTransposed:")
for i = 1, cols do
    local row = {}
    for j = 1, rows do
        row[j] = tablex.array2d_get(transposed, i, j)
    end
    print(tablex.pretty(row))
end

-- Permutations and combinations
print("\nPermutations and Combinations:")
local items = {"a", "b", "c"}
print("Permutations:", tablex.pretty(tablex.permute(items)))
print("Combinations (2):", tablex.pretty(tablex.combinations(items, 2)))

-- Pretty printing
print("\nPretty Printing:")
local complex = {
    name = "test",
    numbers = {1, 2, 3},
    nested = {
        a = "hello",
        b = "world",
        c = {x = 1, y = 2}
    }
}
print("Complex structure:", tablex.pretty(complex, 2))

-- Table loading
print("\nTable Loading:")
local data = tablex.load("{a = 1, b = 2, c = {d = 3}}")
print("Loaded table:", tablex.pretty(data)) 