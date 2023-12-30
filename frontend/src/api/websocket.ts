import { HOST } from "./constants";

var socket = new WebSocket("ws://" + HOST + "/ws")
type Callback = (msg: MessageEvent<string>) => void

export type VALID_MESSAGE_TYPE = "COUNTRY" | "CONTEXT" | "CHAT"
export const [MESSAGE_TYPE_COUNTRY, MESSAGE_TYPE_CONTEXT, MESSAGE_TYPE_CHAT]: VALID_MESSAGE_TYPE[] = ["COUNTRY", "CONTEXT", "CHAT"]

export class Websocket {

    // Listen for new messages to display
    public static connect(cb: Callback): void {
        console.log("connecting...")

        socket.onopen = () => {
            console.log("Successfully Connected")
        }

        socket.onmessage = msg => {
            // console.log(msg)
            cb(msg)
        }

        socket.onclose = event => {
            console.log("Socket Closed Connection:", event)
        }
    }

    // Send to backend
    public static sendMsg(msg: string, msgType: VALID_MESSAGE_TYPE): void {
        console.log("sending msg: ", msg);
        socket.send(Websocket.msgToJson(msg, msgType))
    }

    private static msgToJson(msg: string, msgType: VALID_MESSAGE_TYPE): string {
        const json = JSON.stringify({content: msg, contentType: msgType})
        console.log(json)
        return json
    }
}