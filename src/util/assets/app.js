$(function() {
  var proxies;
  $.ajax({
    method: "POST",
    url: "/proxies",
  })
  .done(function(msg) {
    proxies = $.parseJSON(msg);

    for (var i = 0; i < proxies.length; i++) {
      (function(idx){
        $.ajax({
          method: "POST",
          url: "/api/" + proxies[i],
        })
        .done(function(msg) {
          $('#browser').jsonbrowser(msg, {'collapsed': true});
          // $('#events').html('');
          // prepare divs for events
          $('#events').append('<div class="' + proxies[idx] + '"></div>');
          $('#navs ul').append('<li role="presentation"><a href="#">' + proxies[idx] + '</a></li>');
          $('#navs ul li:last-of-type').click(function(e) {
            $('#navs ul li').removeClass('active');
            $(this).addClass('active');
            $('#events div').not('.' + proxies[idx]).hide();
            $('#events div').not('.' + proxies[idx]).removeClass('active');
            $('#events div' + '.' + proxies[idx]).show();
          });
          getEvents(proxies[idx], $.parseJSON(msg).EventsList);
          if (i == proxies.length) {
            $('#events div').hide();
            $('#events div:first-of-type').show();
            $('#navs ul li:first-of-type').addClass('active');
          }
        });
      })(i);
    }
  });

  var getEvents = function(name, events) {
    var tpl = $("#tpl-events").html();
    var table = $("#tpl-events-table").html();
    $('#events .' + name).append($(_.template(table)()))
    for(var i= 0; i < events.length; i++) {
      var event = events[i];
      (function(ev){
        var $a = $('<a href="#" class=\'btn btn-default\'>Show details</a>')
        var $ev = $(_.template(tpl)({"Event":ev}))
        $a.click(function(e) {
          e.preventDefault();
          $('#request xmp').html(atob(ev['Req']));
          $('#response .headers pre').html(atob(ev['Resp']));
          //$('#response .body pre').html(atob(event['RespBody']));
          try {
            var json = atob(ev['RespBody'])
            JSON.parse(json) // trick for exception
            $('#response .body').jsonbrowser(json, {'collapsed': true});
          } catch (e) {
            $('#response .body').html('<xmp>' + atob(ev['RespBody']) + '</xmp>');
          }
          $("#myModal").modal('show')
        })
        $ev.find('td:last').append($a);
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
