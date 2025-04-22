import path from "path";
import fs from "fs";

export async function GET(req: Request) {
  const logFilePath = path.join(process.cwd(), '../logs', 'shoreline.log');

  // SSE headers
  req.headers.set('Content-Type', 'text/event-stream');
  req.headers.set('Cache-Control', 'no-cache');
  req.headers.set('Connection', 'keep-alive');

  // The response is a readable stream
  const stream = new ReadableStream({
    start(controller) {
      // Function to send data via SSE
      const sendMessage = (data: string) => {
        controller.enqueue(`data: ${data.replaceAll(/\n/g, '$NEWLINE$')}\n\n`);
      };

      // Step 1: Send the entire current content of the log file
      const fileContent = fs.readFileSync(logFilePath, 'utf8');
      sendMessage(fileContent); // Send the current log file content immediately

      // Step 2: Start watching the log file for new lines (tail -f functionality)
      let fileSize = fs.statSync(logFilePath).size;

      // Watch for file changes
      const fileWatcher = fs.watch(logFilePath, { encoding: 'utf8' }, (eventType) => {
        if (eventType === 'change') {
          // When the file changes, read new data from the file
          const newFileSize = fs.statSync(logFilePath).size;

          // If file size increased, read the new part of the log file
          if (newFileSize > fileSize) {
            const newData = fs.createReadStream(logFilePath, {
              encoding: 'utf8',
              start: fileSize, // Start reading from the point where we last left off
              end: newFileSize,
            });

            newData.on('data', (chunk) => {
              sendMessage(chunk); // Send the new data (log entries) to the client
            });

            fileSize = newFileSize; // Update the file size tracker
          }
        }
      });

      // Cleanup: Close the file watcher and stream when the client disconnects
      req.signal.addEventListener('abort', () => {
        fileWatcher.close(); // Stop watching the file
        controller.close();   // Close the stream
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