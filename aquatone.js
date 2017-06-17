var fs        = require('fs');
var Nightmare = require('nightmare');
var nightmare = Nightmare({
  width: 1024,
  height: 768,
  switches: {
    'ignore-certificate-errors': true
  }
});

function debug(message) {
  if (env["AQUATONE_DEBUG"]) {
    console.log(message);
  }
}

var USER_AGENTS = [
  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
  "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0",
  "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
  "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0",
  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.96 Safari/537.36",
  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.1 Safari/603.1.30",
  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
  "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0",
  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
  "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.96 Safari/537.36",
  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
  "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:53.0) Gecko/20100101 Firefox/53.0",
  "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
  "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
  "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
  "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
];

var HTTP_RESPONSES = {
  100: "Continue",
  101: "Switching Protocols",
  200: "OK",
  201: "Created",
  202: "Accepted",
  203: "Non-Authoritative Information",
  204: "No Content",
  205: "Reset Content",
  206: "Partial Content",
  300: "Multiple Choices",
  301: "Moved Permanently",
  302: "Found",
  303: "See Other",
  304: "Not Modified",
  305: "Use Proxy",
  306: "(Unused)",
  307: "Temporary Redirect",
  400: "Bad Request",
  401: "Unauthorized",
  402: "Payment Required",
  403: "Forbidden",
  404: "Not Found",
  405: "Method Not Allowed",
  406: "Not Acceptable",
  407: "Proxy Authentication Required",
  408: "Request Timeout",
  409: "Conflict",
  410: "Gone",
  411: "Length Required",
  412: "Precondition Failed",
  413: "Request Entity Too Large",
  414: "Request-URI Too Long",
  415: "Unsupported Media Type",
  416: "Requested Range Not Satisfiable",
  417: "Expectation Failed",
  500: "Internal Server Error",
  501: "Not Implemented",
  502: "Bad Gateway",
  503: "Service Unavailable",
  504: "Gateway Timeout",
  505: "HTTP Version Not Supported"
};

function randomIpAddress() {
  var result = [];
  for (var i = 0; i < 4; i++) {
    result[i] = Math.floor(Math.random() * (255 - 1 + 1)) + 1;
  }
  return result.join(".");
}

function prettyHeader(header) {
  var words = header.split("-");
  for (var i = 0; i < words.length; i++) {
    words[i] = words[i].charAt(0).toUpperCase() + words[i].slice(1);
  }
  return words.join("-");
}

var url                   = process.argv[2];
var vhost                 = process.argv[3];
var htmlDestination       = process.argv[4];
var headersDestination    = process.argv[5];
var screenshotDestination = process.argv[6];
var visit = {
  success: true,
  url: url,
  vhost: vhost
};

nightmare
  .useragent(USER_AGENTS[Math.floor(Math.random() * USER_AGENTS.length)])
  .goto(url, {
    'Host': vhost,
    'Accept-Language': 'en-US,en;q=0.8',
    'Accept-Encoding': '',
    'Via': '1.1 ' + randomIpAddress(),
    'X-Forwarded-For': randomIpAddress()
  })
  .then(function(result) {
    visit.code   = result.code;
    visit.status = result.code + " " + (HTTP_RESPONSES[visit.code] || "(unknown)")
    var headers  = {};
    for (var k in result.headers) {
      headers[prettyHeader(k)] = result.headers[k].join(" ");
    }
    visit.headers = headers;
    return nightmare
      .wait(1000)
      .evaluate(function() {
        document.body.bgColor = 'white';
      })
      .html(htmlDestination, 'HTMLOnly')
      .then(function(result) {
        return nightmare
          .screenshot(screenshotDestination, {
            x: 0,
            y: 0,
            width: 1024,
            height: 768
          })
          .end()
          .then(function() {
            headers = [];
            for (var k in visit.headers) {
              headers.push(k + ": " + visit.headers[k]);
            }
            fs.writeFile(headersDestination, headers.join("\r\n") + "\r\n", function(err) {
              console.log(JSON.stringify(visit));
              process.exit(0);
            });
          })
      })
  })
  .catch(function(error) {
    console.log(JSON.stringify({
      success: false,
      url: url,
      vhost: vhost,
      error: error.message,
      code: error.code,
      details: error.details
    }));
    process.exit(1);
  });
