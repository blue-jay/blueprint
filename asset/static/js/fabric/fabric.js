requirejs(['js/fabric/api', 'axios'], function (api, axios) {

    $(document).on('click', "input.pushPabric", function(e){
        // var apiFabric = new api.dataAPI("http://192.168.56.101:80/");
        // var data = JSON.stringify({
        //     ext: 1,
        //     type: 2,
        // });
        // var opt = {
        //     data: data,
        // };
        //
        // apiFabric.createFabric("fabric", opt);
        //
        axios.post('/fabric',{
            ext: 1,
            type: 2,
        }).then(function (response) {
            console.log(response);
        }).catch(function (error) {
            console.log(error);
        });
    });

});