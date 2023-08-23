import React, { useEffect, useState } from 'react';

function App() {
    const [messages, setMessages] = useState<any[]>([]);
    const [ws, setWs] = useState<WebSocket | null>(null);

    const connect = () => {
        if (ws) {
            console.log("WebSocket is already connected.");
            return;
        }

        // construct the WebSocket URL based on the current window location
        const url = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/ws`;

        // initialize the WebSocket connection
        const newWs = new WebSocket(url + "?userId=jleva");

        newWs.addEventListener('open', function (event) {
            console.log('WebSocket is open now.', event);
        });

        newWs.addEventListener('message', function (event) {
            console.log('Message from server:', event.data);
            try {
                const messageJson = JSON.parse(event.data);
                setMessages((prevMessages) => [...prevMessages, messageJson]);
            } catch (e) {
                console.error('Could not parse message as JSON:', e);
            }
        });

        setWs(newWs);
    }

    const disconnect = () => {
        if (ws) {
            ws.close();
            setWs(null);
            console.log("WebSocket is closed now.");
        } else {
            console.log("WebSocket is already closed.");
        }
    }

    useEffect(() => {
        return () => {
            if (ws) {
                ws.close();
            }
        };
    }, [ws]);

    return (
        <>
            <button onClick={connect}>Connect</button>
            <button onClick={disconnect}>Disconnect</button>

            {
                messages.map((m, index) => (
                    <pre key={index}>
                        {JSON.stringify(m, null, 2)}
                    </pre>
                ))
            }
        </>
    );
}

export default App;
