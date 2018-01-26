function deepFreeze(o) {
    Object.freeze(o);
    Object.getOwnPropertyNames(o).forEach(function (prop) {
        if (o.hasOwnProperty(prop)
            && o[prop] !== null
            && (typeof o[prop] === "object" || typeof o[prop] === "function")
            && !Object.isFrozen(o[prop])) {
            deepFreeze(o[prop]);
        }
    });
    return o;
}

(function () {
    var deepFreeze = function (o) {
        Object.freeze(o);
        Object.getOwnPropertyNames(o).forEach(function (prop) {
            if (o.hasOwnProperty(prop)
                && o[prop] !== null
                && (typeof o[prop] === "object" || typeof o[prop] === "function")
                && !Object.isFrozen(o[prop])) {
                deepFreeze(o[prop]);
            }
        });
        return o;
    };

    var newConstant = function (parent, property, value) {
        Object.defineProperty(parent, property, {
            writable: false,
            enumerable: true,
            configurable: false,
            value: value
        });
    };

    var newProperty = function (parent, property, value, getter, setter) {
        var configuration = {
            writable: (setter !== null),
            enumerable: true,
            configurable: false
        };
        if (value) configuration.value = value;
        if (getter) {
            configuration.get = function () {
                return getter();
            };
        }
        if (setter) {
            configuration.set = function (newValue) {
                setter(newValue);
            };
        }
        Object.defineProperty(parent, property, configuration);
    };

    var unprivileged = function (call, id, callerId, scoped_system_register) {
        var System = {};
        System.register = function (name, deps, declare) {
            scoped_system_register(name, deps, declare, id, callerId);
        };
        return call();
    };

    var privileged = function () {

    };

    return {
        deepFreeze: deepFreeze,
        newConstant: newConstant,
        newProperty: newProperty,
        unprivileged: unprivileged
    };
})();
