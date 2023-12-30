import { HOST } from "./constants";

var socket = new WebSocket("ws://" + HOST + "/ws")
type Callback = (msg: MessageEvent<string>) => void


export class Websocket {

    // Listen for new messages to display
    public static connect(cb: Callback): void {
        console.log("connecting...")

        socket.onopen = () => {
            console.log("Successfully Connected")
        }

        socket.onmessage = msg => {
            console.log(msg)
            cb(msg)
        }

        socket.onclose = event => {
            console.log("Socket Closed Connection:", event)
        }
    }

    // Send to backend
    public static sendMsg(msg: string): void {
        console.log("sending msg: ", msg);
        socket.send(msg)
    }
}