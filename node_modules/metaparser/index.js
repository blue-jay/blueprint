/*jslint node:true, nomen:true*/

(function () {

    'use strict';

    var fs = require('fs'),
        path = require('path'),
        cheerio = require('cheerio'),
        _ = require('underscore'),
        mkdirp = require('mkdirp'),
        async = require('async');

    module.exports = function (params) {

        var options = _.defaults(params || {}, {
            source: null,
            add: null,
            remove: null,
            out: null,
            callback: null
        });

        async.waterfall([
            function (callback) {
              var decodeData = function (error, data) {
                var $ = cheerio.load(data, {
                    decodeEntities: false
                });
                callback(error, $);
              }

              if (options.source) {
                fs.readFile(options.source, decodeData);
              }
              else {
                decodeData(undefined, options.data);
              }
            },
            function ($, callback) {
                if (options.remove) {
                    _.each(typeof options.remove === 'string' ? [options.remove] : options.remove, function (code) {
                        $(code).remove();
                    });
                    callback(null, $);
                } else {
                    callback(null, $);
                }
            },
            function ($, callback) {
                if (options.add) {
                    var target = $('head').length > 0 ? $('head') : $.root();
                    _.each(typeof options.add === 'string' ? [options.add] : options.add, function (code) {
                        target.append(code);
                    });
                    callback(null, $);
                } else {
                    callback(null, $);
                }
            },
            function ($, callback) {
                if (options.out) {
                    mkdirp(path.dirname(options.out), function (error) {
                        if (error) {
                            callback(error, $);
                        } else {
                            fs.writeFile(options.out, $.html(), function (error) {
                                callback(error, $);
                            });
                        }
                    });
                } else {
                    callback(null, $);
                }
            }
        ], function (error, $) {
            if (options.callback) {
                return options.callback(error, $.html().replace(/^\s*[\r\n]/gm, ''));
            }
            if (error) {
                throw error;
            }
            return;
        });

    };

}());
