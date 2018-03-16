var metaparser = require('./index.js');
var fs = require('fs');

// Data loaded from a file
metaparser({
    source: 'test/index.html',
    add: '<link rel="author" href="humans.txt" />',
    remove: 'meta[name="author"]',
    out: 'test/index-metaparser.html',
    callback: function (error, data) {
        console.log(error);
        console.log(data);
    }
});

// Data are transmitted directly
var data = fs.readFileSync('test/index.html');
metaparser({
    data: data,
    add: '<link rel="author" href="humans.txt" />',
    remove: 'meta[name="author"]',
    out: 'test/index-metaparser.html',
    callback: function (error, data) {
        console.log(error);
        console.log(data);
    }
});
