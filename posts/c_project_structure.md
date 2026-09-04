# C Project Structure

/ -> root of project, tool configs go here (.clangd, .gitignore, .clang-format, etc)
/examples/ -> playground + can have ci to test whether all compile
/build/ -> anything build system related (CMakeLists.txt, scripts, configs, compile_commands.json)
/dist/ -> output of build system
/src/ -> project source code
/src/\<module\> -> a whole module that results in a _single_ executable or dll
/src/\<module\>/\*\*/\*\_test.c -> unit test
/tools/ -> tools written in c that are used in the project (either in the build system or by the application itself)

## My Own Build System

TODO
