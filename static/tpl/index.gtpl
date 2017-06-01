{{define "head"}}
<style>
ul {
    width: 50%;
    margin: auto;
}
</style>
{{end}}
{{define "title"}}{{.Title}} {{end}}
{{define "body"}}
{{.Body}}
<form class="form-horizontal" id="register-data" action="/" method="post">
<fieldset>
<div class="jumbotron text-center">
  <h1>SBERBANK TELEX KEY GENERATOR</h1>
  <p>Утилита формирования ключа для TELEX</p> 
</div>
<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="suminput">Укажите сумму сделки</label>  
  <div class="col-md-4">	
	<div class="input-group">	  
	  <span class="input-group-addon">Руб</span>
	  <input type="text" class="form-control" aria-label="Окргулено до рубля" id="suminput" name="suminput" type="text" placeholder="10000000">
	  <span class="input-group-addon">.00</span>
	</div>
	<span class="help-block">Сумма сделки</span>    
  </div>
</div>


<!-- Select Basic -->
<div class="form-group">
  <label class="col-md-4 control-label" for="selectbasic">Валюта</label>
  <div class="col-md-2">
    <select id="selectbasic" name="selectbasic" class="form-control">
      <option value="RUB">RUB</option>
      <option value="NOC">NOC</option>
      <option value="OTH">OTH</option>
      <option value="USD">USD</option>
      <option value="EUR">EUR</option>
    </select>
  </div>
</div>

<!-- Textarea -->
<div class="form-group">
  <label class="col-md-4 control-label" for="textarea">Сопровождающее сообщение TELEX</label>
  <div class="col-md-4">                     
    <textarea class="form-control" id="textarea" name="textarea" rows="6" cols="60">ТЕКСТ СООБЩЕНИЯ
</textarea>
  </div>
</div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="textinput">Дата сделки</label>  
  <div class="col-md-4">
  <input id="dateinput" name="dateinput" type="text" placeholder="{{ .DateNow }}" value="{{ .DateNow }}" class="form-control input-md" readonly>
  <span class="help-block">Дата сделки</span>  
  </div>
</div>

<!-- Button -->
<div class="form-group">
  <label class="col-md-4 control-label" for="savebutton">Сформировать сообщение</label>
  <div class="col-md-4">
    	<button id="savebutton" name="savebutton" class="btn btn-primary">Сформировать</button>   
  </div>
</div>
<!-- Button -->
<div class="form-group">
  <label class="col-md-4 control-label" for="exitbutton">Завершение приложения</label>
  <div class="col-md-4">
	<button id="exitbutton" name="exitbutton" class="btn btn-danger">Выход</button>    
  </div>
</div>



</fieldset>
</form>
{{end}}
{{define "scripts"}}
<script type="text/javascript" language="javascript">
$('#savebutton').click(function () {
$('#savebutton').prop('disabled', true);
var data = $("#register-data").serializeObject();
data["Post"]="SaveButton"
alert(JSON.stringify(data));
$.ajax({                 /* start ajax function to send data */
        url: "/",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',
        error: function () { alert("POST Handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {
			//alert("REG: "+data);
			// handle AJAX redirection
			if (JSON.parse(data) == "SaveOk") {
				alert("Сообщение сформировано успешно.");
				window.location = '/success';
			}
			if (JSON.parse(data) == "SaveNotOk"){
				alert("Данные введены с ошибкой. Получатель не был добавлен.");
				$('#savebutton').prop('disabled', false);
			}
						
        }
    }); 
});
$('#exitbutton').click(function () {
$('#exitbutton').prop('disabled', true);
var data = {};
data["Post"]="ExitButton"
alert(JSON.stringify(data));
$.ajax({                 /* start ajax function to send data */
        url: "/",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',
        error: function () { alert("POST Handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {
			//alert("REG: "+data);
			// handle AJAX redirection
			if (JSON.parse(data) == "ExitOk") {
				alert("Приложение остановлено, работа с приложением завершена.");
				window.location = 'about:blank';
			}
						
        }
    }); 
});


$.fn.serializeObject = function()
{
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
            if (!o[this.name].push) {
                o[this.name] = [o[this.name]];
            }
            o[this.name].push(this.value || '');
        } else {
            o[this.name] = this.value || '';
        }
    });
    return o;
};

</script>
{{end}}