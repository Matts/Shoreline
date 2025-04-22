export async function GET(req: Request) {
  const stream = new ReadableStream({
    start(controller) {
      // Set the SSE headers
      req.headers.set('Content-Type', 'text/event-stream');
      req.headers.set('Cache-Control', 'no-cache');
      req.headers.set('Connection', 'keep-alive');

      // Function to send message as SSE
      const sendMessage = (data: any) => {
        controller.enqueue(`data: ${JSON.stringify(data)}\n\n`);
      };

      // Send message every 5 seconds
      const intervalId = setInterval(() => {
        sendMessage({ message: 'New log entry', timestamp: new Date().toISOString() });
      }, 5000);

      // Cleanup when client disconnects
      req.signal.addEventListener('abort', () => {
        clearInterval(intervalId); // Stop sending data
        controller.close(); // Close the stream
      });
    },
  });

  // Return the Response with stream
  return new Response(stream, {
    headers: {
      'Content-Type': 'text/event-stream',
    },
  });
}