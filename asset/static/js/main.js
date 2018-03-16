requirejs.config({
    baseUrl: 'static',
    paths : {
        jquery: 'js/jquery.min',
        // axios: 'js/axios/axios.min',
        bootstrap: 'js/bootstrap.min',
        all: 'js/all.min',
        underscore: 'js/underscore-min',
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