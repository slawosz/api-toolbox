$(function() {
  var proxies;
  $.ajax({
    method: "POST",
    url: "/proxies",
    //dataType: "json",
    //data: JSON.stringify(data),
  })
  .done(function(msg) {
    proxies = $.parseJSON(msg);

    for (var i = 0; i < proxies.length; i++) {
      (function(idx){
        $.ajax({
          method: "POST",
          url: "/api/" + proxies[i],
          //dataType: "json",
          //data: JSON.stringify(data),
        })
        .done(function(msg) {
          $('#browser').jsonbrowser(msg, {'collapsed': true});
          // $('#events').html('');
          // prepare divs for events
          $('#events').append('<div class="' + proxies[idx] + '"></div>');
          $('#navs ul').append('<li role="presentation"><a href="#">' + proxies[idx] + '</a></li>');
          getEvents(proxies[idx], $.parseJSON(msg).EventsList);
        });
      })(i);
    }
  });

  var getEvents = function(name, events) {
    console.log(events);
    var tpl = $("#tpl-events").html();
    var table = $("#tpl-events-table").html();
    $('#events .' + name).append($(_.template(table)()))
    console.log(events.length);
    for(var i= 0; i < events.length; i++) {
      var event = events[i];
      console.log(event);
      (function(ev){
        console.log(event);
        var $a = $('<a href="#" class=\'btn btn-default\'>Show details</a>')
        var $ev = $(_.template(tpl)({"Event":ev}))
        $a.click(function(e) {
          e.preventDefault();
          $('#request pre').html(atob(ev['Req']));
          $('#response .headers pre').html(atob(ev['Resp']));
          //$('#response .body pre').html(atob(event['RespBody']));
          var json = atob(ev['RespBody'])
          try {
            $('#response .body').jsonbrowser(json, {'collapsed': true});
          } catch (e) {
            $('#response .body').html('<pre>' + json + '</pre>');
          }
          $("#myModal").modal('show')
        })
        $ev.find('td:last').append($a);
        console.log($('#events .' + name + ' tbody'));
        $('#events .' + name + ' tbody').append($ev);
      })(event)
    }
  }

  var helpers = {
    fdate: function(date) {
      var d = new Date(date)
      return d.toISOString()
    }
  }

  $('#collapse-all').on('click', function(e) {
    e.preventDefault();
    $.jsonbrowser.collapseAll('#browser');
  });

  $('#expand-all').on('click', function(e) {
    e.preventDefault();
    $.jsonbrowser.expandAll('#browser');
  });

  $('#search').on('keyup', function(e) {
    e.preventDefault();
    $.jsonbrowser.search('#browser', $(this).val());
  });
  $('#search').focus().trigger('keyUp');
});
