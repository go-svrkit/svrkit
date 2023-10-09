# svrkit

### Guidelines

Please follow the repository guidelines and standards mentioned below.


### Making Changes
---

#### Code

- Please use the directory structure of the repository.
- Use meaningful variable names.
- Strictly use snake_case (underscore_separated)  in filenames.
- If you have added or modified code, please make sure the code compiles before submitting.


#### Documentation

- Make sure you put useful comments in your code. Do not comment on obvious things.
- If you have modified/added documentation, please ensure that your language is concise and must not contain grammatical errors.
- Do not update [`README.md`] along with other changes.


#### Test

- Make sure to add examples and test cases in your `filename_test.go` file.
- If you find an algorithm or document without tests, please feel free to create a pull request or issue describing suggested changes.
- Please try to add one or more `Test` functions that will invoke the algorithm implementation on random test data with the expected output.

### Benchmark
---

- Make sure to add examples and benchmark cases in your `filename_test.go` or `filename_bench.go` if you want separated test and benchmark files.
- Please try to add one or more `Benchmark` functions that will invoke the algorithm implementation.
- For running the benchmark, you could use this command `go test -bench=.` for more details, read this article [Using Subtests and Sub-benchmarks](https://go.dev/blog/subtests)



#### New File Name guidelines

- Use lowercase words without ``"_"`` for the file name
- Use ``"_"`` as a separator only for `_test.go` or `_bench.go`
- For instance

```markdown
MyNewGoFile.GO       is incorrect
my_new_go_file.go    is incorrect
mynewgofile.go       is the correct format
mynewgofile_test.go  is the correct format
```

- It will be used to dynamically create a directory of files and implementation.
- File name validation will run on Docker to ensure validity.
- Check out `Go` [Package names](https://go.dev/blog/package-names) roles

#### New Directory guidelines

- We recommend adding files to existing directories as much as possible.
- Use lowercase words with ``"_"`` as separator ( no spaces or ```"-"``` allowed )
- For instance

```markdown
SomeNew Fancy-Category          is incorrect
some_new_fancy_category         is correct
```

- Filepaths will be used to dynamically create a directory of our algorithms.
- Filepath validation will run on GitHub Actions to ensure compliance.

#### Commit Guidelines

- It is recommended to keep your changes grouped logically within individual commits. It's easier to understand changes that are logically spilled across multiple commits. Try to modify just one or two files in the same directory.

```bash
git add filexyz.go
git commit -m "your message"
```

Examples of commit messages with semantic prefixes:

```markdown
fix: XYZ algorithm bug
feat: add XYZ algorithm
test: add test for XYZ algorithm
docs: add comments and explanation to XYZ algorithm
```

Common prefixes:

- fix: A bug fix
- feat: A new feature
- docs: Documentation changes
- test: Correct existing tests or add new ones

