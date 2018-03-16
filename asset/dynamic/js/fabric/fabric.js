define([
    "jquery",
    "fabric/api"], function ($, fabricApi) {
    var api = fabricApi.dataAPI("/login");
    var data = JSON.stringify({
        ext: 1,
        type: 2,
    })
    var opt = {
        data: data,
    }
    api.createFabric("fabric", opt);
});