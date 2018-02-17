tsVersion(ts.version);
function transpiler(source) {
    var result = ts.transpileModule(source, {
        compilerOptions: {
            moduleResolution: "node",
            module: "system",
            target: "es5",
            tsconfig: false,
            noImplicitAny: false,
            alwaysStrict: true,
            inlineSourceMap: true,
            diagnostics: true,
            strictPropertyInitialization: true,
            allowJs: false,
            downlevelIteration: true,
            noLib: true,
            typeRoots: [
                "scripts/types"
            ],
            lib: [
                "lib/libbase.d.ts"
            ]
        }
    });

    return result.outputText;
}