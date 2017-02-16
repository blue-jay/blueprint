// Modules
var gulp        = require('gulp');
var favicon     = require('gulp-real-favicon');
var fs          = require('fs');
var runSequence = require('run-sequence'); // Using until gulp v4 is released
var reload      = require('gulp-livereload');
var sync        = require('gulp-sync')(gulp).sync;
var child       = require('child_process');
var util        = require('gulp-util');
var path        = require('path');
var os          = require('os');

// Enviroment variables
var env = JSON.parse(fs.readFileSync('./env.json'))
var folderAsset = env.Asset.Folder;
var folderView = env.View.Folder;

// Other variables
var faviconData = folderAsset + '/dynamic/favicon/data.json';

// Application server
var server = null;

// SASS Task
gulp.task('sass', function() {
	var sass = require('gulp-sass');
	var ext = require('gulp-ext-replace');
	gulp.src(folderAsset + '/dynamic/sass/**/*.scss')
		// Available for outputStyle: expanded, nested, compact, compressed
		.pipe(sass({outputStyle: 'expanded'}).on('error', sass.logError))
		.pipe(gulp.dest(folderAsset + '/static/css/'));
	return gulp.src(folderAsset + '/dynamic/sass/**/*.scss')
		// Available for outputStyle: expanded, nested, compact, compressed
		.pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError))
		.pipe(ext('.min.css'))
		.pipe(gulp.dest(folderAsset + '/static/css/'))
		.pipe(reload());
});

// JavaScript Task
gulp.task('javascript', function() {
	var concat = require('gulp-concat');
	var minify = require('gulp-minify');
	var babel = require('gulp-babel');
	return gulp.src(folderAsset + '/dynamic/js/*.js')
		.pipe(babel({
			presets: ['es2015']
		}))
		.pipe(concat('all.js'))
		.pipe(minify({
			ext:{
				src:'.js',
				min:'.min.js'
			}
		}))
		.pipe(gulp.dest(folderAsset + '/static/js/'))
		.pipe(reload());
});

// jQuery Task
gulp.task('jquery', function() {
	return gulp.src('node_modules/jquery/dist/jquery.min.*')
		.pipe(gulp.dest(folderAsset + '/static/js/'));
});

// Bootstrap Task
gulp.task('bootstrap', function() {
	gulp.src('node_modules/bootstrap/dist/css/bootstrap-theme.min.*')
		.pipe(gulp.dest(folderAsset + '/static/css/'));
	gulp.src('node_modules/bootstrap/dist/css/bootstrap.min.*')
		.pipe(gulp.dest(folderAsset + '/static/css/'));
	gulp.src('node_modules/bootstrap/dist/fonts/*')
		.pipe(gulp.dest(folderAsset + '/static/fonts/'));
	return gulp.src('node_modules/bootstrap/dist/js/bootstrap.min.js')
		.pipe(gulp.dest(folderAsset + '/static/js/'));
});

// Underscore Task
gulp.task('underscore', function() {
	return gulp.src('node_modules/underscore/underscore-min.*')
		.pipe(gulp.dest(folderAsset + '/static/js/'));
});

// Favicon Generation and Injection Task
gulp.task('favicon', function() {
	runSequence('favicon-generate', 'favicon-inject');
});

// Generate the icons. This task takes a few seconds to complete.
// You should run it at least once to create the icons. Then,
// you should run it whenever RealFaviconGenerator updates its
// package (see the favicon-update task below).
gulp.task('favicon-generate', function(done) {
	var favColor = '#525252';
	favicon.generateFavicon({
		masterPicture: folderAsset + '/dynamic/favicon/logo.png',
		dest: folderAsset + '/static/favicon/',
		iconsPath: '/static/favicon/',
		design: {
			ios: {
				pictureAspect: 'backgroundAndMargin',
				backgroundColor: favColor,
				margin: '14%'
			},
			desktopBrowser: {},
			windows: {
				pictureAspect: 'noChange',
				backgroundColor: favColor,
				onConflict: 'override'
			},
			androidChrome: {
				pictureAspect: 'noChange',
				themeColor: favColor,
				manifest: {
					name: 'Blueprint',
					display: 'browser',
					orientation: 'notSet',
					onConflict: 'override',
					declared: true
				}
			},
			safariPinnedTab: {
				pictureAspect: 'silhouette',
				themeColor: favColor
			}
		},
		settings: {
			scalingAlgorithm: 'Mitchell',
			errorOnImageTooSmall: false
		},
		versioning: {
			paramName: 'v1.0',
			paramValue: '3eepn6WlLO'
		},
		markupFile: faviconData
	}, function() {
		done();
	});
});

// Inject the favicon markups in your HTML pages. You should run
// this task whenever you modify a page. You can keep this task
// as is or refactor your existing HTML pipeline.
gulp.task('favicon-inject', function() {
	return gulp.src([folderView + '/partial/favicon.tmpl'])
		.pipe(favicon.injectFaviconMarkups(JSON.parse(fs.readFileSync(faviconData)).favicon.html_code))
		.pipe(gulp.dest(folderView + '/partial/'));
});

// Check for updates on RealFaviconGenerator (think: Apple has just
// released a new Touch icon along with the latest version of iOS).
// Run this task from time to time. Ideally, make it part of your
// continuous integration system.
gulp.task('favicon-update', function(done) {
	var currentVersion = JSON.parse(fs.readFileSync(faviconData)).version;
	return favicon.checkForUpdates(currentVersion, function(err) {
		if (err) {
			throw err;
		}
	});
});

// Monitor Go files for changes
gulp.task('server:watch', function() {
	// Restart application
	gulp.watch([
		'*/**/*.tmpl',
		'env.json'
	], ['server:spawn']);
	
	// Rebuild and restart application server
	gulp.watch([
		'*.go',
		'*/**/*.go'
	], [
		'server:build',
		'server:spawn'
	]);
});

// Build application from source
gulp.task('server:build', function() {
	var build = child.spawnSync('go', ['build']);
	if (build.stderr.length) {
		var lines = build.stderr.toString()
		.split('\n').filter(function(line) {
		return line.length
		});
		for (var l in lines)
			util.log(util.colors.red(
			'Error (go build): ' + lines[l]
		));
		notifier.notify({
			title: 'Error (go build)',
			message: lines
		});
	}
	return build;
});

// Spawn an application process
gulp.task('server:spawn', function() {
	if (server)
		server.kill();
	
	// Get the application name based on the folder
	var appname = path.basename(__dirname);

	// Spawn application server
	if (os.platform() == 'win32') {
		server = child.spawn(appname + '.exe');
	} else {
		server = child.spawn('./' + appname);
	}
	
	// Trigger reload upon server start
	server.stdout.once('data', function() {
		reload.reload('/');
	});
	
	// Pretty print server log output
	server.stdout.on('data', function(data) {
		var lines = data.toString().split('\n')
		for (var l in lines)
		if (lines[l].length)
		util.log(lines[l]);
	});
	
	// Print errors to stdout
	server.stderr.on('data', function(data) {
		process.stdout.write(data.toString());
	});
});

// Main watch function.
gulp.task('watch', ['server:build'], function() {
	// Start the listener (use with the LiveReload Chrome Extension)
	reload.listen();

	// Watch the assets
	gulp.watch(folderAsset + '/dynamic/sass/**/*.scss', ['sass']);
	gulp.watch(folderAsset + '/dynamic/js/*.js', ['javascript']);
	
	return gulp.start(sync([
		'server:watch',
		'server:spawn'
	]));
});

// Init - every task
gulp.task('init', ['sass', 'javascript', 'jquery', 'bootstrap', 'underscore', 'favicon', 'server:build']);

// Default - only run the tasks that change often
gulp.task('default', ['sass', 'javascript', 'server:build']);
