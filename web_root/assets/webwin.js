(function(e,t) {
    window.addEventListener('load', init);
    var project_name="WebWin";
    var selected;
    function init() {
        M.AutoInit();
        //Add an load some internal modules.
        load_module("internal", "main", "main");
        //Add menu buttons from plugins
        $.post("/api/plugins?action=list", function(data, status) {
            var obj = jQuery.parseJSON(data);
            jQuery.each( obj["plugins"], function( i, val ) {
                if (val["hidden"])
                    return;
                var el=$("<li class=\"bold\"><a class=\"waves-effect waves-teal\"><i class=\"material-icons\">"+val["icon"]+"</i>"+val["name"]+"</a></li>")
                    .attr("module_id",val["id"])
                    .on("click", function() {
                        button_module(this);
                    });
                $("#menu").append(el);
            });
        });
    }
    $.fn.ignore = function(sel) {
        return this.clone().find(sel || ">*").remove().end();
      };
    async function load_module(directory, module_id, module_name) {
        return await $.get(directory+"/"+module_id+"/index.html", function(data, status){
            $('title').text(module_name.charAt(0).toUpperCase() + module_name.slice(1) +" - "+project_name);
            $('#page_name').text(module_name.charAt(0).toUpperCase() + module_name.slice(1));
            $('main #module_content').empty().append(data);
            result=true;
        }).fail(function() {
            M.toast({text: "Failed to load plugin "+module_id+"!"});
        });
    }
    async function button_module(element) {
        var module_id=jQuery(element).attr("module_id");
        // Code from here: https://stackoverflow.com/questions/11347779
        var module_name=jQuery(element).children("a").contents().filter(function() {return this.nodeType == 3;}).text();
        if (await load_module("plugins", module_id, module_name)) {
            console.log("Loaded module: \""+module_id+"\".");
            jQuery(selected).removeClass("active");
            selected=element;
            jQuery(element).addClass("active")
        }
    }
})();
