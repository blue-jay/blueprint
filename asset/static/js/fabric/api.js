define([
"jquery",
], function($){
    var dataAPI = function(base_url){
        this.base_url = base_url;
    };
    var default_options = {
        type: 'GET',
        contentType: 'application/JSON',
        cache: false,
        async: true,
        processData: false,
        success: null,

    };

    var update = function(d1, d2){
        $.map(d2, function(i, key){
            d1[key] = d2[key];
        });
        return d1;
    };

    var ajax_dafaults = function(options){
        var d = {};
        update(d, default_options || {});
        update(d, options || {});
        return d;
    };

    dataAPI.prototype.api_request = function(path, options){
        options = options || {};
        options = ajax_dafaults(options);
        url = this.base_url + path;
        $.ajax(url, options);
    };

    dataAPI.prototype.getStaticData = function(path, options){
        var res;
        var opt = {
            dataType: 'json',
            async: false,
            success: function(data){
                // console.log(data);
                res = data;
                // console.log('get data succeed.');
                // console.log(res);
            },
            error: function(XMLHttpRequest, textStatus, errorThrown){
                console.log(textStatus, errorThrown, XMLHttpRequest.status, XMLHttpRequest.readyStatus);
                res = JSON.parse('[]');
                console.log('get data file failed.');
            },
        };
        update(opt, options || {});
        // console.log('before ajax.');
        this.api_request(path, opt);
        // console.log('after ajax');
        return res;
    };

    dataAPI.prototype.createFabric = function(path, options){

        var opt = {
            type: 'POST',
            dataType: 'json',
            contentType: "application/json",
            async: false,
            processData: false,
            success: function(data, status, jqXHR){
                console.log(data, status, jqXHR);
            },
            error: function(XMLHttpRequest, textStatus, errorThrown){
                console.log(textStatus, errorThrown, XMLHttpRequest.status, XMLHttpRequest.readyStatus);
            },
        };
        update(opt, options || {});
        this.api_request(path, opt);
    };
    return {'dataAPI': dataAPI};

});
