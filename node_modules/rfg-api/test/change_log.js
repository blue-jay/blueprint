/*jslint node:true*/

'use strict';

var rfg = require('../index.js').init();
var assert = require('assert');

describe('Change log', function() {
  describe('#changeLog()', function() {
    it("should return changes since a certain versions", function(done) {
      rfg.changeLog("0.9", function(err, versions) {
        assert.equal(err, undefined);

        assert(versions.length > 1);
        assert.equal(versions[0].version, '0.10');

        done();
      });
    });

    it("should return all changes", function(done) {
      rfg.changeLog(undefined, function(err, versions) {
        assert.equal(err, undefined);

        assert(versions.length > 10);
        assert.equal(versions[0].version, '0.1');

        done();
      });
    });

  });
});
