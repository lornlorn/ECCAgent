$(function () {
    // 调试exec接口
    $('#exec_submit').click(function () {
        var params = {};
        params['from'] = 'test';
        params['data'] = {};
        params['data']['envs'] = [];
        params['data']['args'] = [];
        var dict = {};
        dict['env1'] = $('#env1').val();
        dict['env2'] = $('#env2').val();
        params['data']['envs'].push(dict['env1']);
        params['data']['envs'].push(dict['env2']);
        params['data']['cmd'] = $('#cmd').val();
        dict['arg1'] = $('#arg1').val();
        dict['arg2'] = $('#arg2').val();
        params['data']['args'].push(dict['arg1']);
        params['data']['args'].push(dict['arg2']);

        console.log('REQUEST : ' + JSON.stringify(params));

        $.ajax({
            // url: '/test/ajax',
            url: '/exec',
            type: 'POST',
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify(params),
            async: 'true',
            dataType: 'json',
            success: function (result) {
                console.log('RESPONSE : ' + JSON.stringify(result));
                console.log("请求成功");
                // alert('成功');
                // window.close();
                $('#ret').val(JSON.stringify(result));
            },
            error: function (result) {
                console.log("请求失败");
            },
            complete: function () {
                console.log("Ajax finish");
            },
        });
    });

    // 新增cron
    $('#CronAdd_submit').click(function () {
        var params = {};
        params['from'] = 'test';
        params['data'] = {};
        params['data']['CronEnvs'] = [];
        params['data']['CronArgs'] = [];

        params['data']['CronName'] = $('#CronName').val();
        params['data']['CronSpec'] = $('#CronSpec').val();

        var dict = {};
        dict['CronEnv1'] = $('#CronEnv1').val();
        dict['CronEnv2'] = $('#CronEnv2').val();
        params['data']['CronEnvs'].push(dict['CronEnv1']);
        params['data']['CronEnvs'].push(dict['CronEnv2']);
        params['data']['CronCmd'] = $('#CronCmd').val();
        dict['CronArg1'] = $('#CronArg1').val();
        dict['CronArg2'] = $('#CronArg2').val();
        params['data']['CronArgs'].push(dict['CronArg1']);
        params['data']['CronArgs'].push(dict['CronArg2']);

        params['data']['CronStatus'] = $('#CronStatus').val();
        params['data']['CronDesc'] = $('#CronDesc').val();

        console.log('REQUEST : ' + JSON.stringify(params));

        $.ajax({
            // url: '/test/ajax',
            url: '/cron',
            type: 'POST',
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify(params),
            async: 'true',
            dataType: 'json',
            success: function (result) {
                console.log('RESPONSE : ' + JSON.stringify(result));
                console.log("请求成功");
                // alert('成功');
                // window.close();
                $('#CronRet').val(JSON.stringify(result));
            },
            error: function (result) {
                console.log("请求失败");
            },
            complete: function () {
                console.log("Ajax finish");
            },
        });
    });

    // 删除cron
    $('#CronDelete_submit').click(function () {
        var params = {};
        params['from'] = 'test';
        params['data'] = {};
        params['data']['CronName'] = $('#CronName').val();
        params['data']['CronUuid'] = $('#CronUuid').val();

        console.log('REQUEST : ' + JSON.stringify(params));

        $.ajax({
            // url: '/test/ajax',
            url: '/cron',
            type: 'DELETE',
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify(params),
            async: 'true',
            dataType: 'json',
            success: function (result) {
                console.log('RESPONSE : ' + JSON.stringify(result));
                console.log("请求成功");
                // alert('成功');
                // window.close();
                $('#CronRet').val(JSON.stringify(result));
            },
            error: function (result) {
                console.log("请求失败");
            },
            complete: function () {
                console.log("Ajax finish");
            },
        });
    });

    // 更新cron
    $('#CronUpdate_submit').click(function () {
        var params = {};
        params['from'] = 'test';
        params['data'] = {};
        params['data']['CronEnvs'] = [];
        params['data']['CronArgs'] = [];

        params['data']['CronName'] = $('#CronName').val();
        params['data']['CronSpec'] = $('#CronSpec').val();

        var dict = {};
        dict['CronEnv1'] = $('#CronEnv1').val();
        dict['CronEnv2'] = $('#CronEnv2').val();
        params['data']['CronEnvs'].push(dict['CronEnv1']);
        params['data']['CronEnvs'].push(dict['CronEnv2']);
        params['data']['CronCmd'] = $('#CronCmd').val();
        dict['CronArg1'] = $('#CronArg1').val();
        dict['CronArg2'] = $('#CronArg2').val();
        params['data']['CronArgs'].push(dict['CronArg1']);
        params['data']['CronArgs'].push(dict['CronArg2']);

        params['data']['CronStatus'] = $('#CronStatus').val();
        params['data']['CronDesc'] = $('#CronDesc').val();
        params['data']['CronUuid'] = $('#CronUuid').val();

        console.log('REQUEST : ' + JSON.stringify(params));

        $.ajax({
            // url: '/test/ajax',
            url: '/cron',
            type: 'PUT',
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify(params),
            async: 'true',
            dataType: 'json',
            success: function (result) {
                console.log('RESPONSE : ' + JSON.stringify(result));
                console.log("请求成功");
                // alert('成功');
                // window.close();
                $('#CronRet').val(JSON.stringify(result));
            },
            error: function (result) {
                console.log("请求失败");
            },
            complete: function () {
                console.log("Ajax finish");
            },
        });
    });

});