import {native__hello_world} from "native";

export class TestPlugin {

    public helloWorld(name: String): String {
        return native__hello_world(name);
    }
}
