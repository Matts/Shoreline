"use client";

import {useEffect, useState} from 'react';

const LogStream = () => {
  const [logs, setLogs] = useState<string>();
  const [showingBadge, setShowingBadge] = useState(false);
  const [originalFaviconHref, setOriginalFaviconHref] = useState<string | null>(null);

  useEffect(() => {
    const eventSource = new EventSource('/api/sse');

    // Handle incoming messages
    eventSource.onmessage = (event) => {
      let data = event.data;
      data = data.replaceAll(/\$NEWLINE\$/g, '\n'); // Replace custom newline marker with actual newline
      // event data is just string append
      setLogs((prev) => (prev ? prev + data : data));

      createFaviconWithBadge();
    };

    // Handle errors
    eventSource.onerror = (error) => {
      console.error('Error with SSE connection:', error);
      eventSource.close(); // Close the connection on error
    };

    // Cleanup on unmount
    return () => {
      eventSource.close();
    };
  }, []);


    // Function to create a favicon with a red dot
    const createFaviconWithBadge = () => {

      const originalFavicon = document.querySelector('link[rel="icon"]');


      if (!originalFavicon) {
        console.error('No favicon found');
        return;
      }

      if (!originalFaviconHref) {
        setOriginalFaviconHref(originalFavicon.href);
      }

      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');

      canvas.width = 32;
      canvas.height = 32;

      // Draw the original favicon (a simple colored background as a placeholder)
      // ctx.fillStyle = 'blue';  // Placeholder color, use your actual favicon here
      // ctx.fillRect(0, 0, canvas.width, canvas.height);

      // Draw a red dot badge on top-right corner
      ctx.beginPath();
      ctx.arc(24, 8, 6, 0, 2 * Math.PI);  // Position the red dot
      ctx.fillStyle = 'red';
      ctx.fill();

      // Create the favicon image from the canvas
      const faviconUrl = canvas.toDataURL();
      document.querySelector('link[rel="icon"]').href = faviconUrl;
      setShowingBadge(true);
    };

  const resetFavicon = () => {
    document.querySelector('link[rel="icon"]').href = originalFaviconHref;
    setShowingBadge(false);
  };

    useEffect(() => {
      if(!showingBadge) {
        return;
      }

      const resetTimeout = setTimeout(() => {
        resetFavicon();
      }, 2000); // Reset after 5 seconds
      return () => {
        clearTimeout(resetTimeout);
      };
    }, [showingBadge])



  return (
    <div>
      <h3>Logs:</h3>
      <span className={"whitespace-pre-wrap"}>{logs}</span>
    </div>
  );
};

export default LogStream;