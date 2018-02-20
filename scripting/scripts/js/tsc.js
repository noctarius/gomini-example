tsVersion(ts.version);

function transpiler(source) {
    var result = ts.transpileModule(source, {
        compilerOptions: {
            moduleResolution: "node",
            module: "System",
            target: "es5",
            isolatedModules: true,
            importHelpers: true,
            tsconfig: false,
            noImplicitAny: false,
            alwaysStrict: true,
            inlineSourceMap: true,
            diagnostics: true,
            strictPropertyInitialization: true,
            allowJs: false,
            downlevelIteration: true,
            noLib: true,
            declaration: true,
            typeRoots: [
                "scripts/types"
            ],
            lib: [
                "lib/libbase.d.ts"
            ]
        },
        reportDiagnostics: true,
        transformers: []
    });

    for (var i = 0; i < result.diagnostics.length; i++) {
        ts.sys.write(result.diagnostics[i]);
    }

    return result.outputText;
}
