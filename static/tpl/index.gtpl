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
	  <input type="text" class="form-control" aria-label="Окргулено до рубля" id="suminput" name="suminput" type="text" placeholder="0">
	  <span class="input-group-addon">.00</span>
	</div>
	<span class="help-block">Сумма сделки</span>    
  </div>
</div>


<!-- Select Basic -->
<div class="form-group">
  <label class="col-md-4 control-label" for="selectcur">Валюта</label>
  <div class="col-md-2">
    <select id="selectcur" name="selectcur" class="form-control">
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
    <textarea class="form-control" id="textarea" name="textarea" rows="6" cols="60">

    TEST CODE FROM YOUR MESSAGE WAS SUCCESSFULLY CHECKED    

    

</textarea>
  </div>
</div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="dateinput">Дата сделки</label>  
  <div class="col-md-4">
  <input id="dateinput" name="dateinput" type="text" placeholder="{{ .DateNow }}" value="{{ .DateNow }}" class="form-control input-md" readonly>
  <span class="help-block">Дата сделки</span>  
  </div>
</div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="seqcounter">Номер данного сообщения</label>  
  <div class="col-md-4">
  <input id="seqcounter" name="seqcounter" type="text" placeholder="" value="{{ .SeqCnt }}" class="form-control input-md" readonly>
  <span class="help-block">Номер данного сообщения (в году)</span>  
  </div>
</div>

<!-- Text input-->
<div class="form-group">
  <label class="col-md-4 control-label" for="telexkey">Ключ TELEX</label>  
  <div class="col-md-4">
  <input id="telexkey" name="telexkey" type="text" placeholder="" value="КЛЮЧ НЕ БЫЛ ВЫЧИСЛЕН" class="form-control input-md" readonly>
  <span class="help-block">Номер данного сообщения (в году)</span>  
  </div>
</div>

<!-- Textarea -->
<div class="form-group">
  <label class="col-md-4 control-label" for="calclog">Рассчет</label>
  <div class="col-md-4">                     
    <textarea class="form-control" id="calclog" name="calclog" rows="6" cols="60">
</textarea>
  </div>
</div>


<!-- Button -->
<div class="form-group">
  <label class="col-md-4 control-label" for="calcbutton">Рассчитать ключ</label>
  <div class="col-md-4">
    	<button id="calcbutton" name="calcbutton" class="btn btn-primary">Рассчитать</button>   
  </div>
</div>

<!-- Button -->
<div class="form-group">
  <label class="col-md-4 control-label" for="savebutton">Сформировать сообщение</label>
  <div class="col-md-4">
    	<button id="savebutton" name="savebutton" class="btn btn-success">Сформировать</button>   
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
//alert(JSON.stringify(data));
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
			if (JSON.parse(data) == "SaveNotOkSUM"){
				alert("Сумма сделки введена с ошибкой.");
				$('#savebutton').prop('disabled', false);
			}
			if (JSON.parse(data) == "SaveNotOk"){
				alert("Данные введены с ошибкой.");
				$('#savebutton').prop('disabled', false);
			}

						
        }
    }); 
});
$('#exitbutton').click(function () {
$('#exitbutton').prop('disabled', true);
var data = {};
data["Post"]="ExitButton"
//alert(JSON.stringify(data));
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

$('#calcbutton').click(function () {
$('#calcbutton').prop('disabled', true);
var data = $("#register-data").serializeObject();
data["Post"]="CalcButton"
//alert(JSON.stringify(data));
$.ajax({                 /* start ajax function to send data */
        url: "/",
        type: 'POST',
        datatype: 'json',
        contentType: 'application/json; charset=UTF-8',
        error: function () { alert("POST Handshake didn't go through") }, /* call disconnect function */
        data: JSON.stringify(data),
        success: function (data) {			
			arr = JSON.parse(data);
			//alert("Ключ TELEX: "+arr[1]);
			// handle AJAX redirection
			if (arr[0] == "CalcOk") {				
				document.getElementById("telexkey").value = arr[1];
				document.getElementById("calclog").value = arr[2];
				
			}
			if (JSON.parse(data) == "CalcNotOk"){
				alert("Данные введены с ошибкой.");
				$('#calcbutton').prop('disabled', false);
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