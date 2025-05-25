-- Create a directory structure
os.execute("mkdir -p test_dir/subdir")
os.execute("echo 'Hello, World!' > test_dir/file1.txt")
os.execute("echo 'Another file' > test_dir/subdir/file2.txt")

-- Create a TAR archive
tar.create("archive.tar", "test_dir")
print("Created archive.tar")

-- Create a compressed TAR archive
tar.create("archive.tar.gz", "test_dir", true)
print("Created archive.tar.gz")

-- List contents of the archive
local files = tar.list("archive.tar")
print("\nArchive contents:")
for i, file in ipairs(files) do
    print(string.format("%s (%d bytes, %s)", file.name, file.size, file.type))
end

-- Extract the archive to a new directory
tar.extract("archive.tar", "extracted")
print("\nExtracted archive to 'extracted' directory")

-- List contents of the compressed archive
local compressed_files = tar.list("archive.tar.gz")
print("\nCompressed archive contents:")
for i, file in ipairs(compressed_files) do
    print(string.format("%s (%d bytes, %s)", file.name, file.size, file.type))
end

-- Extract the compressed archive
tar.extract("archive.tar.gz", "extracted_gz")
print("\nExtracted compressed archive to 'extracted_gz' directory")