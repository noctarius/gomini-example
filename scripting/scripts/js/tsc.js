function transpiler(source) {
    var result = ts.transpileModule(source, {
        compilerOptions: {
            module: "System",
            target: "es5",
            tsconfig: false,
            noImplicitAny: false,
            alwaysStrict: true,
            inlineSourceMap: true,
            typeRoots: [
                "scripts/types"
            ]
        }
    });

    return result.outputText;
}