package html

var SayHelloHtml = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go HTTP Server</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container { 
            background: white;
            padding: 40px;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            max-width: 600px;
            width: 100%;
            animation: slideUp 0.6s ease-out;
        }
        @keyframes slideUp {
            from { opacity: 0; transform: translateY(30px); }
            to { opacity: 1; transform: translateY(0); }
        }
        h1 { 
            color: #333;
            margin-bottom: 20px;
            font-size: 2em;
        }
        h1::before {
            content: '🎉 ';
            animation: bounce 2s infinite;
        }
        @keyframes bounce {
            0%, 100% { transform: translateY(0); }
            50% { transform: translateY(-10px); }
        }
        p { 
            color: #555;
            line-height: 1.6;
            margin-bottom: 20px;
        }
        .info { 
            background: linear-gradient(135deg, #e8f4fd 0%, #d4e7f7 100%);
            padding: 20px;
            border-radius: 8px;
            margin: 25px 0;
            border-left: 4px solid #667eea;
        }
        .info strong {
            display: block;
            color: #333;
            margin-bottom: 12px;
            font-size: 1.1em;
        }
        .info ul {
            list-style: none;
            padding: 0;
        }
        .info li {
            padding: 8px 0;
            padding-left: 24px;
            position: relative;
            color: #444;
        }
        .info li::before {
            content: '✓';
            position: absolute;
            left: 0;
            color: #667eea;
            font-weight: bold;
        }
        .timestamp { 
            color: #667eea;
            font-weight: 600;
            font-size: 1em;
            background: #f0f4ff;
            padding: 2px 8px;
            border-radius: 4px;
        }
        .footer {
            margin-top: 25px;
            padding-top: 20px;
            border-top: 2px solid #f0f0f0;
            font-style: italic;
            color: #666;
        }
        .badge {
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 0.85em;
            font-weight: 600;
            margin-left: 8px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>HTTP Server Working!</h1>
        <p>Congratulations! You've successfully built a TCP server that speaks HTTP. <span class="badge">Go</span></p>
        
        <div class="info">
            <strong>What just happened:</strong>
            <ul>
                <li>Your browser sent an HTTP request over TCP</li>
                <li>Our server parsed the raw TCP data</li>
                <li>We generated a proper HTTP response</li>
                <li>The browser rendered this HTML!</li>
            </ul>
        </div>
        
        <p>This response was generated at: <span class="timestamp" id="timestamp">2025-10-24 14:30:45</span></p>
        
        <div class="footer">
            <p>Next: We'll learn to parse HTTP requests and create a proper router!</p>
        </div>
    </div>
    
    <script>
        // Update timestamp every second to show it's dynamic
        setInterval(() => {
            const now = new Date();
            const formatted = now.getFullYear() + '-' + 
                String(now.getMonth() + 1).padStart(2, '0') + '-' + 
                String(now.getDate()).padStart(2, '0') + ' ' +
                String(now.getHours()).padStart(2, '0') + ':' + 
                String(now.getMinutes()).padStart(2, '0') + ':' + 
                String(now.getSeconds()).padStart(2, '0');
            document.getElementById('timestamp').textContent = formatted;
        }, 1000);
    </script>
</body>
</html>
`

var BadRequestHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>400 Bad Request</title>
    <style>
        body { 
            font-family: Arial, sans-serif;
            background: grey;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0;
            padding: 20px;
        }
        .container { 
            background: white;
            padding: 40px;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            max-width: 500px;
            text-align: center;
        }
        .code { 
            font-size: 3em;
            font-weight: bold;
            color: #ff6b6b;
            margin-bottom: 10px;
        }
        h1 { 
            color: #333;
            font-size: 1.5em;
            margin-bottom: 15px;
        }
        p { 
            color: #666;
            line-height: 1.5;
        }
        .timestamp { 
            color: #ff6b6b;
            font-weight: 600;
            background: fff0f0;
            padding: 2px 8px;
            border-radius: 4px;
            display: inline-block;
            margin-top: 20px;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="code">400</div>
        <h1>⚠️ Bad Request</h1>
        <p>The server couldn't understand your request. Please check your HTTP formatting and try again.</p>
        <div class="timestamp" id="timestamp">2025-10-24 14:30:45</div>
    </div>
</body>
</html>
`
