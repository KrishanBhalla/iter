var socket = new WebSocket("ws://localhost:3000/ws")


type Callback = (msg: MessageEvent<string>) => void

// Listen for new messages to display
let connect = (cb: Callback) => {
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
let sendMsg = (msg: string) => {
    console.log("sending msg: ", msg);
    socket.send(msg)
}

export { connect, sendMsg }