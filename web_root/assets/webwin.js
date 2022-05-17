(function(e,t) {
    window.addEventListener('load', init);
    var project_name="WebWin";
    var selected;
    function init() {
        M.AutoInit();
        load_module("000_main", "Main page");
        $('#logo-container').on("click", function() {
            jQuery(selected).removeClass("active");
            load_module("000_main", "Main page");
        });
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
    async function load_module(module_id, module_name) {
        return await $.get("plugins/"+module_id+"/index.html", function(data, status){
            $('title').attr("module_id",module_id).text(module_name.charAt(0).toUpperCase() + module_name.slice(1) +" - "+project_name);
            $('#page_name').text(module_name.charAt(0).toUpperCase() + module_name.slice(1));
            $('main #module_content').empty().append(data);
            console.log("Loaded module: \""+module_id+"\".");
            result=true;
        }).fail(function() {
            M.toast({text: "Failed to load plugin "+module_id+"!"});
        });
    }
    async function button_module(element) {
        var module_id=jQuery(element).attr("module_id");
        // Code from here: https://stackoverflow.com/questions/11347779
        var module_name=jQuery(element).children("a").contents().filter(function() {return this.nodeType == 3;}).text();
        if (await load_module(module_id, module_name)) {
            jQuery(selected).removeClass("active");
            selected=element;
            jQuery(element).addClass("active")
        }
    }
})();
