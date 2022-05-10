(function(e,t) {
    window.addEventListener('load', init);
    var project_name="WebWin";
    var selected;
    function init() {
        M.AutoInit();
        $.post("/api/plugins", function(data, status) {
            console.log(data);
            var obj = jQuery.parseJSON(data);
            jQuery.each( obj["plugins"], function( i, val ) {
                var el=$("<li class=\"bold\"><a class=\"waves-effect waves-teal\">"+val+"</a></li>")
                    .on("click", function() {
                        selected=this;
                        button_module(this);
                    });
                $("#menu").append(el);
            });
        });
        load_module("internal", "main");
    }
    function load_module(directory, module_name) {
        $.get(directory+"/"+module_name+"/index.html", function(data, status){
            $('title').text(module_name.charAt(0).toUpperCase() + module_name.slice(1) +" "+project_name);
            $('#page_name').text(module_name.charAt(0).toUpperCase() + module_name.slice(1));
            $('main #module_content').empty().append(data);
            return true;
        }).fail(function() {
            M.toast({text: "Failed to load plugin "+module_name+"!"});
            return false;
        });
    }
    function button_module(element) {
        var module_name=jQuery(element).children("a").text();
        if (load_module("plugins", module_name)) {
            jQuery(element).removeClass("active");
            selected=element;
            jQuery(element).addClass("active")
        }
    }
})();
