var project_name="WebWin";

var current_module_id;
var current_module_name;

async function sleep(time) {
    return new Promise(resolve => setTimeout(resolve, time));
}

function ToastIcon(icon, text) {
    M.toast({unsafeHTML: "<span class=\"material-icons\">"+icon+"</span>"+ text});
}

function ToastError(text) {
    ToastIcon("error", "Error (see console): "+text);
}

async function actionGetStatus(action_id) {
    var tmp;
    await $.post("/api/plugins?action=get_status&id="+action_id, function(data, status) {
        tmp=data;
    }).fail(function(jqXHR, textStatus, error) {
        ToastError(error);
    });
    return tmp;
}

async function actionRequest(module_id, action, param) {
    var tmp;
    await $.post("/api/plugins/"+module_id+"?action="+action+"&param="+param, function(data, status) {
        tmp=data;
    }).fail(function(jqXHR, textStatus, error) {
        ToastError(error);
    });
    return tmp;
}

async function actionRequestAndWait(module_id, action, param) {
    var id = await actionRequest(module_id, action, param);
    var result = {"status": "pending"};
    while (result.status == "pending") {
        result = await actionGetStatus(id);
        await sleep(400);
    }
    return result;
}

async function clearJobs() {
    return await $.post("/api/plugins?action=jobs_clear", function(data, status){
        console.log("Clear jobs success: " + data);
    }).fail(function(jqXHR, textStatus, error) {
        ToastError(error);
    });
}

(function(e,t) {
    window.addEventListener('load', init);
    var selected;
    // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Proxy
    // Pretty fast way to get url params
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
      });
    function init() {
        M.AutoInit();
        if (params.id==null) {
            load_module("000_main", "Main page");
        } else {
            selected=params.id;
            load_module(params.id, params.name);
        }
        $('#logo-container').on("click", function() {
            jQuery(selected).removeClass("active");
            load_module("000_main", "Main page");
        });
        //Add menu buttons from plugins
        $.post("/api/plugins?action=list", function(data, status) {
            //In fact, server returns string with json, so we need to parse it
            var obj = jQuery.parseJSON(data);
            jQuery.each( obj["plugins"], function( i, val ) {
                if (val["hidden"])
                    return;
                var el=$("<li class=\"bold\"><a class=\"waves-effect waves-teal\"><i class=\"material-icons\">"+val["icon"]+"</i>"+val["name"]+"</a></li>")
                    .attr("module_id",val["id"])
                    .on("click", function() {
                        button_module(this);
                    });
                $('#menu').append(el);
            });
        });
    }
    async function load_module(module_id, module_name) {
        return await $.get("plugins/"+module_id+"/index.html", function(data, status){
            $('title').attr("module_id",module_id).text(module_name.charAt(0).toUpperCase() + module_name.slice(1) +" - "+project_name);
            $('#page_name').text(module_name.charAt(0).toUpperCase() + module_name.slice(1));
            $('main #module_content').empty().append(data);
            current_module_id=module_id;
            current_module_name=module_name;
            window.history.pushState({}, module_name, "?id="+module_id+"&name="+module_name);
            console.log("Loaded module: \""+module_id+"\".");
        }).fail(function() {
            M.toast({text: "Failed to load plugin "+module_id+"!"});
        });
    }
    async function button_module(element) {
        if (selected==element)
            return;
        var module_id=jQuery(element).attr("module_id");
        // Code from here: https://stackoverflow.com/questions/11347779
        // We ignore any other elements except text
        var module_name=jQuery(element).children("a").contents().filter(function() {return this.nodeType == 3;}).text();
        if (await load_module(module_id, module_name)) {
            jQuery(selected).removeClass("active");
            selected=element;
            jQuery(element).addClass("active");
        }
    }
})();
