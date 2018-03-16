requirejs.config({
    baseUrl: 'static',
    paths : {
        jquery: 'js/jquery.min',
        bootstrap: 'js/bootstrap.min',
        all: 'js/all.min',
        underscore: 'js/underscore-min',
        // axios: "js/axios/axios"
    },
    // shim:{
    //     axios:{
    //         deps:[],
    //         exports: 'axios',
    //     }
    // },
});

// define(["jquery"], function($) {
//     $(function() {
//     });
// });