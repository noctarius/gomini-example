(function () {
    var pushIfAbsent = function (properties, property) {
        for (var el in properties) {
            if (el === property) {
                return false;
            }
        }
        properties.push(property);
        return true;
    };

    var collectProperties = function (object) {
        var properties = [];
        var obj = object;
        while (true) {
            if (obj === null || obj === Object.prototype) {
                break;
            }
            Object.keys(obj).forEach(function (e) {
                if (pushIfAbsent(properties, e)) {
                }
            });
            obj = Object.getPrototypeOf(obj);
        }
        return properties;
    };

    var debug = function (properties) {
        Object.getOwnPropertyNames(properties).forEach(function (obj) {
            var props = properties[obj];
            console.log(obj + ": " + props);
        })
    };

    return function (object, adaptNull, adaptFunction, adaptArray, adaptObject, adaptProperty) {
        var properties = collectProperties(object);
        //debug(properties);
        for (var i = 0; i < properties.length; i++) {
            var property = properties[i];
            var prop = object[property];
            if (prop === null) {
                adaptNull(property);
            } else if (typeof prop === 'function') {
                var constructor = function () {
                    var args = Array.prototype.slice.call(arguments);
                    args.unshift(null);
                    return new (Function.prototype.bind.apply(prop, args));
                };
                adaptFunction(property, prop, prop, constructor);
            } else if (Array.isArray(prop)) {
                adaptArray(property, prop);
            } else if (typeof prop === 'object') {
                adaptObject(property, prop);
            } else {
                var descriptor = Object.getOwnPropertyDescriptor(object, property);
                adaptProperty(property, descriptor);
            }
        }
    }
})();