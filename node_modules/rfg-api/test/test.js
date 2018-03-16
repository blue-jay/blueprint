/*jslint node:true*/

'use strict';

var api = require('../index.js').init();
var request = require('./request.json');
var rimraf = require('rimraf');

var assert = require('assert');
var path = require('path');
var rfg = require('../index.js').init();
var fs = require('fs');

describe('RFG Api', function() {

  beforeEach(function() {
    if (! fs.existsSync(path.join(__dirname, 'output'))) {
      fs.mkdirSync(path.join(__dirname, 'output'));
    }
  });

  afterEach(function() {
    rimraf.sync(path.join(__dirname, 'output'));
  });

  describe('#fileToBase64()', function() {
    it('should return the content of a file encoded in base64', function(done) {
      rfg.fileToBase64(path.join(__dirname, 'input', 'very_small.png'), function(error, base64) {
        if (error) throw error;
        assert.equal(
          'iVBORw0KGgoAAAANSUhEUgAAAAIAAAADCAIAAAA2iEnWAAAACXBIWXMAAAsTAAALEwEAmpwYAAAA' +
          'B3RJTUUH3woWBxkR5IGL1wAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUH' +
          'AAAAHElEQVQI1wXBgQAAAACDsAiCuD9TLN9IXbhSUuJAYwXpQ37pHAAAAABJRU5ErkJggg==',
          base64);
        done();
      });
    });

    it('should return an error when the file does not exist', function(done) {
      rfg.fileToBase64('oops', function(error, base64) {
        assert.notEqual(error, undefined);
        done();
      });
    });
  });

  describe('#fileToBase64Sync()', function() {
    it('should return the content of a file encoded in base64', function() {
      assert.equal(
        rfg.fileToBase64Sync(path.join(__dirname, 'input', 'very_small.png')),
        'iVBORw0KGgoAAAANSUhEUgAAAAIAAAADCAIAAAA2iEnWAAAACXBIWXMAAAsTAAALEwEAmpwYAAAA' +
        'B3RJTUUH3woWBxkR5IGL1wAAAB1pVFh0Q29tbWVudAAAAAAAQ3JlYXRlZCB3aXRoIEdJTVBkLmUH' +
        'AAAAHElEQVQI1wXBgQAAAACDsAiCuD9TLN9IXbhSUuJAYwXpQ37pHAAAAABJRU5ErkJggg==');
    });
  });

  describe('#generateFavicon()', function() {
    this.timeout(15000);

    it('should generate a favicon', function(done) {
      rfg.fileToBase64(path.join(__dirname, 'input', 'master_picture.png'), function(error, base64) {
        assert.equal(error, undefined);
        var req = {
          "api_key": "f26d432783a1856427f32ed8793e1d457cc120f1",
          "master_picture": {
            "type": "inline",
            "content": base64
          },
          "files_location": {
            "type": "path",
            "path": "favicons/"
          },
          "favicon_design": {
            "ios": {
              "picture_aspect": "background_and_margin",
              "margin": "4",
              "background_color": "#123456"
            }
          },
          "settings": {
            "compression": 1,
            "scaling_algorithm": "NearestNeighbor"
          }
        };
        rfg.generateFavicon(req, path.join(__dirname, 'output'), function(err, result) {
          assert.equal(err, undefined);

          // Make sure iOS icons were generated, but not desktop icons
          assert(fs.statSync(path.join(__dirname, 'output', 'apple-touch-icon.png')).isFile());
          assert(! fs.existsSync(path.join(__dirname, 'output', 'favicon.ico')));

          // Make sure some code is returned
          assert(result.favicon.html_code);
          assert(result.favicon.html_code.length > 500);
          assert(result.favicon.html_code.length < 1500);

          done();
        });
      });
    });

    it('should generate a favicon based on an SVG image', function(done) {
      rfg.fileToBase64(path.join(__dirname, 'input', 'master_picture.svg'), function(error, base64) {
        assert.equal(error, undefined);
        var req = {
          "api_key": "f26d432783a1856427f32ed8793e1d457cc120f1",
          "master_picture": {
            "type": "inline",
            "content": base64
          },
          "files_location": {
            "type": "path",
            "path": "favicons/"
          },
          "favicon_design": {
            "desktop_browser": {}
          }
        };
        rfg.generateFavicon(req, path.join(__dirname, 'output'), function(err, result) {
          assert.equal(err, undefined);

          // Make sure desktop icons were generated, but not iOS icons
          assert(! fs.existsSync(path.join(__dirname, 'output', 'apple-touch-icon.png')));
          assert(fs.statSync(path.join(__dirname, 'output', 'favicon.ico')).isFile());

          // Make sure some code is returned
          assert(result.favicon.html_code);
          assert(result.favicon.html_code.length > 200);
          assert(result.favicon.html_code.length < 1000);

          done();
        });
      });
    });
  });

  describe('#injectFaviconMarkups()', function() {
    it('should inject favicon code', function(done) {
      var markups = [
        '<link rel="icon" type="image/png" href="favicons/favicon-192x192.png" sizes="192x192">',
        '<link rel="icon" type="image/png" href="favicons/favicon-160x160.png" sizes="160x160">'
      ];
      var fileContent = fs.readFileSync(path.join(__dirname, 'input', 'test_1.html'));
      rfg.injectFaviconMarkups(fileContent, markups, {}, function(error, html) {
        var expected = fs.readFileSync(path.join(__dirname, 'input', 'test_1_expected_output.html')).toString();
        assert.equal(html, expected);

        done();
      });
    });

    it('should remove existing markups', function(done) {
      var markups = [
        '<link rel="icon" type="image/png" href="favicons/favicon-192x192.png" sizes="192x192">',
        '<link rel="icon" type="image/png" href="favicons/favicon-160x160.png" sizes="160x160">'
      ];
      var fileContent = fs.readFileSync(path.join(__dirname, 'input', 'test_2.html'));
      rfg.injectFaviconMarkups(fileContent, markups, {}, function(error, html) {
        var expected = fs.readFileSync(path.join(__dirname, 'input', 'test_2_expected_output.html')).toString();
        assert.equal(html, expected);

        done();
      });
    });

    it('should inject extra markups', function(done) {
      var markups = [
        '<link rel="icon" type="image/png" href="favicons/favicon-192x192.png" sizes="192x192">',
        '<link rel="icon" type="image/png" href="favicons/favicon-160x160.png" sizes="160x160">'
      ];
      var fileContent = fs.readFileSync(path.join(__dirname, 'input', 'test_2.html'));
      rfg.injectFaviconMarkups(fileContent, markups, {
        add: '<link content="an extra markup">'
      }, function(error, html) {
        var expected = fs.readFileSync(path.join(__dirname, 'input', 'test_2_expected_output_with_extra.html')).toString();
        assert.equal(html, expected);

        done();
      });
    });

    it('should remove extra markups', function(done) {
      var markups = [
        '<link rel="icon" type="image/png" href="favicons/favicon-192x192.png" sizes="192x192">',
        '<link rel="icon" type="image/png" href="favicons/favicon-160x160.png" sizes="160x160">'
      ];
      var fileContent = fs.readFileSync(path.join(__dirname, 'input', 'test_2.html'));
      rfg.injectFaviconMarkups(fileContent, markups, {
        remove: ['meta[name="description"]']
      }, function(error, html) {
        var expected = fs.readFileSync(path.join(__dirname, 'input', 'test_2_expected_output_with_removal.html')).toString();
        assert.equal(html, expected);

        done();
      });
    });

    it('should keep extra markups', function(done) {
      var markups = [
        '<link rel="icon" type="image/png" href="favicons/favicon-192x192.png" sizes="192x192">',
        '<link rel="icon" type="image/png" href="favicons/favicon-160x160.png" sizes="160x160">'
      ];
      var fileContent = fs.readFileSync(path.join(__dirname, 'input', 'test_2.html'));
      rfg.injectFaviconMarkups(fileContent, markups, {
        keep: 'link[rel="icon"]'
      }, function(error, html) {
        var expected = fs.readFileSync(path.join(__dirname, 'input', 'test_2_expected_output_with_keeping.html')).toString();
        assert.equal(html, expected);

        done();
      });
    });
  });

  describe('#camelCaseToUnderscore()', function() {
    it('should turn camel case to underscores', function() {
      // One word
      assert.equal('hello', rfg.camelCaseToUnderscore('hello'));
      // Two words
      assert.equal('hello_world', rfg.camelCaseToUnderscore('helloWorld'));
      // Long string and there are two consecutive uppercase letters
      assert.equal('hello_world_this_is_a_long_string', rfg.camelCaseToUnderscore('helloWorldThisIsALongString'));
      // First letter is uppercased
      assert.equal('hello', rfg.camelCaseToUnderscore('Hello'));
      // No effect on an underscore string
      assert.equal('hello_world', rfg.camelCaseToUnderscore('hello_world'));
      // Numbers
      assert.equal('option1_a', rfg.camelCaseToUnderscore('option1A'));
    });
  });

  describe('#camelCaseToUnderscoreRequest()', function() {
    it('should convert a JS request (camelcase) to an RFG request (underscore)', function() {
      assert.deepEqual(rfg.camelCaseToUnderscoreRequest({}), {});

      assert.equal(rfg.camelCaseToUnderscoreRequest(undefined), undefined);

      assert.deepEqual(rfg.camelCaseToUnderscoreRequest({
        firstEntry: 'firstValue',
        secondEntry: [
          'aValue',
          'anotherValue',
          8,
          {
            aSubHash: 'itsValue',
            scaling_algorithm: 'NearestNeighbor'
          }
        ],
        thirdEntry: {
          firstSubEntry: 'itsValue',
          secondSubEntry: 'anotherValue'
        }
      }), {
        first_entry: 'first_value',
        second_entry: [
          'a_value',
          'another_value',
          8,
          {
            a_sub_hash: 'its_value',
            scaling_algorithm: 'NearestNeighbor'
          }
        ],
        third_entry: {
          first_sub_entry: 'its_value',
          second_sub_entry: 'another_value'
        }
      });
    });
  });
});

describe('Request helpers', function() {
  describe('#escapeJSONSpecialChars()', function() {
    it('should escape special characters', function() {
      assert.equal('hello', rfg.escapeJSONSpecialChars('hello'));
      assert.equal('\"hello\"', rfg.escapeJSONSpecialChars('"hello"'));
      assert.equal('e\&p', rfg.escapeJSONSpecialChars('e&p'));
    });
  });

  describe('#isUrl()', function() {
    it('should set path and URL apart', function() {
      assert( rfg.isUrl('http://www.example.com'));
      assert( rfg.isUrl('https://www.example.com'));
      assert(!rfg.isUrl('/my/project'));
      assert(!rfg.isUrl('images/mu_pic.png'));
    });
  });

  describe('#isBase64()', function() {
    it('should indicate if a string is base64 or not', function() {
      assert( rfg.isBase64('U29tZSByYW5kb20gY29udGVudA=='));
      assert(!rfg.isBase64(path.join(__dirname, 'input', 'small_file.txt')));
    });
  });

  describe('#normalizeMasterPicture()', function() {
    it('should inline file content when necessary', function() {
      assert.deepEqual(rfg.normalizeMasterPicture({
        type: 'inline',
        content: path.join(__dirname, 'input', 'small_file.txt')
      }), {
        type: 'inline',
        content: "U29tZSByYW5kb20gY29udGVudA=="
      });

      assert.deepEqual(rfg.normalizeMasterPicture({
        type: 'inline',
        content: 'U29tZSByYW5kb20gY29udGVudA=='
      }), {
        type: 'inline',
        content: "U29tZSByYW5kb20gY29udGVudA=="
      });

      assert.deepEqual(rfg.normalizeMasterPicture({
        content: path.join(__dirname, 'input', 'small_file.txt')
      }), {
        type: 'inline',
        content: "U29tZSByYW5kb20gY29udGVudA=="
      });

      var urlMP = {
        type: 'url',
        url: 'http://www.example.com/a_picture.png'
      };
      assert.deepEqual(rfg.normalizeMasterPicture(urlMP), urlMP);
    });
  });

  describe('#normalizeAllMasterPictures()', function() {
    it('should inline all master pictures of a request', function() {
      var dummyRequest = {
        master_picture: {
          content: path.join(__dirname, 'input', 'small_file.txt')
        },
        stuff: [
          {
            a: 'b',
            master_picture: {
              type: 'inline',
              content: path.join(__dirname, 'input', 'small_file.txt')
            }
          }
        ]
      };

      var normRequest = {
        master_picture: {
          type: 'inline',
          content: "U29tZSByYW5kb20gY29udGVudA=="
        },
        stuff: [
          {
            a: 'b',
            master_picture: {
              type: 'inline',
              content: "U29tZSByYW5kb20gY29udGVudA=="
            }
          }
        ]
      };

      assert.deepEqual(rfg.normalizeAllMasterPictures(dummyRequest), normRequest);
    });
  });

  describe('#createRequest()', function() {
    it('should generate a RFG API request without settings or versioning', function() {
      assert.deepEqual(rfg.createRequest({
        apiKey: '123azerty',
        masterPicture: path.join(__dirname, 'input', 'small_file.txt'),
        iconsPath: '/path/to/icons',
        design: {
          desktop: {},
          ios: {
            masterPicture: {
              content: path.join(__dirname, 'input', 'small_file.txt'),
            },
            pictureAspect: 'noChange'
          }
        }
      }),{
        api_key: '123azerty',
        favicon_design: {
          desktop: {},
          ios: {
            master_picture: {
              content: "U29tZSByYW5kb20gY29udGVudA==",
              type: 'inline'
            },
            picture_aspect: 'no_change'
          }
        },
        files_location: {
          path: '/path/to/icons',
          type: 'path'
        },
        master_picture: {
          content: "U29tZSByYW5kb20gY29udGVudA==",
          type: 'inline'
        }
      });
    });

    it('should generate a RFG API request with settings or versioning', function() {
      assert.deepEqual(rfg.createRequest({
        apiKey: '123azerty',
        masterPicture: path.join(__dirname, 'input', 'small_file.txt'),
        design: {
          desktop: {}
        },
        settings: {
          compression: 3
        },
        versioning: {
          paramName: 'theName',
          paramValue: '123abc'
        }
      }),{
        api_key: '123azerty',
        favicon_design: {
          desktop: {}
        },
        settings: {
          compression: 3
        },
        versioning: {
          param_name: 'theName',
          param_value: '123abc'
        },
        files_location: {
          type: 'root'
        },
        master_picture: {
          content: "U29tZSByYW5kb20gY29udGVudA==",
          type: 'inline'
        }
      });
    });
  });
});
