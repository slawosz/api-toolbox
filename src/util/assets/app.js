$(function() {
  $.ajax({
    method: "POST",
    url: "/api",
    //dataType: "json",
    //data: JSON.stringify(data),
  })
  .done(function(msg) {
    $('#browser').jsonbrowser(msg, {'collapsed': true});
    $('#events').html('')
    getEvents($.parseJSON(msg))
  });

  var helpers = {
    fdate: function(date) {
      var d = new Date(date)
      return d.toISOString()
    }
  }

  var getEvents = function(events) {
    var tpl = $("#tpl-events").html();
    var table = $("#tpl-events-table").html();
    $('#events').append($(_.template(table)()))
    for(var i= 0; i < events.length; i++) {
      var event = events[i];
      (function(ev){
        var $a = $('<a href="#" class=\'btn btn-default\'>Show details</a>')
        var $ev = $(_.template(tpl)({"Event":ev}))
        $a.click(function(e) {
          e.preventDefault();
          $('#request pre').html(atob(ev['Req']));
          $('#response .headers pre').html(atob(ev['Resp']));
          //$('#response .body pre').html(atob(event['RespBody']));
          var json = atob(ev['RespBody'])
          $('#response .body').jsonbrowser(json, {'collapsed': true});
          $("#myModal").modal('show')
        })
        $ev.find('td:last').append($a)
        $('#events tbody').append($ev)
      })(event)
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
