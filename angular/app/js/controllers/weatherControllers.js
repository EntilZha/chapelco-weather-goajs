var chapelcoWeatherAppControllers = angular.module('chapelcoWeatherApp.controllers');

chapelcoWeatherAppControllers.controller('CurrentWeatherCtrl', ['$scope', '$http', function($scope, $http) {
	$scope.dataLoaded = false;
	$http.get('api/weather/current').success(function(data) {
		$scope.currentWeather = data;
		$scope.dataLoaded = true;
	});
}]);

chapelcoWeatherAppControllers.controller('WeatherChartsCtrl', ['$scope', '$http', function($scope, $http) {
	$http.get('api/weather/past-field-lists/432').success(function(data) {
		$scope.data = data;
	});
	options = {
		xAxis: {
			title: {}
		},
		yAxis: {
			plotLines: [{
				value: 0,
				width: 1,
				color: '#808080'
			}],
			title: {}
		},
		title: {}
	}
	var makeChart = function(title, yTitle, seriesName) {
		options.title = { text: title, x: -20};
		options.xAxis.title.text = 'Time';
		options.xAxis.categories = $scope.data['DATE_TIME'];
		options.xAxis.labels = { rotation: 45, step: 18 };
		options.yAxis.title.text = yTitle;
		options.series = [{
			data: $scope.data[seriesName],
			name: yTitle
		}];
		$('#weather-chart').highcharts(options);
		$('#chart-modal').modal('show');
	};
	$scope.drawChart = function(chartType) {
		if (chartType == "CHN1_DEG") {
			makeChart("Past Temperatures", "Temperature (°C)", chartType);
		} else if (chartType == "CHN1_DEW") {
			makeChart("Past Dew Points", "Temperature (°C)", chartType);
		} else if (chartType == "CHN1_RF") {
			makeChart("Past Relative Humidities", "Relative Humidity (%)", chartType);
		} else if (chartType == "RAIN_SUM") {
			makeChart("Past Rain Sums", "MM Water", chartType);
		} else if (chartType == "PRES_LOC") {
			makeChart("Past Pressure", "hPa", chartType);
		} else if (chartType == "PRES_ABS") {
			makeChart("Past Absolute Pressure", "hPa", chartType);
		}
	};
}]);
