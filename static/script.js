$(document).ready(function() {
  $("#search").on('input', function(event) {
    event.preventDefault();
    animeFilter($(this).val())
  });

  animeFilter = function(query) {
    $("tbody tr").each(function(index, el) {
      if (query == "") {
        $(this).show('slow');
        return;
      }

      var aTitle = $(this).find("#aTitle").text().toLowerCase();
      var aQuery = query.toLowerCase();

      if (aTitle.indexOf(aQuery) > -1) $(this).show('fast');
      else $(this).hide('fast');
    });
  }
});
