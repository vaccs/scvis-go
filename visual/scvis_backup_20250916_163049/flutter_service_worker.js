'use strict';
const MANIFEST = 'flutter-app-manifest';
const TEMP = 'flutter-temp-cache';
const CACHE_NAME = 'flutter-app-cache';

const RESOURCES = {"flutter.js": "888483df48293866f9f41d3d9274a779",
"favicon.png": "5dcef449791fa27946b3d35ad8803796",
"assets/assets/images/playctl_end.png": "00ee42645813d0a205e9d655eee93853",
"assets/assets/images/CstrcpyUnboundExample2.png": "7880a341f4133380b97a7ac85edd536b",
"assets/assets/images/playctl_back.png": "1ddbe4aa3bc016805324d529e3a75efb",
"assets/assets/images/logo.png": "2766030b6b31fdb9845f8006cb09bf30",
"assets/assets/images/playctl_play.png": "3c914893d471a7a2b3de8c44a1026f85",
"assets/assets/images/playctl_start.png": "84ea41401e3abbdd531ca25752bf1070",
"assets/assets/images/CStringFigureSecurity.png": "a60bbbc9d83403045c9895db553814e4",
"assets/assets/images/CmemExample.png": "e4257885a53079ae921466324f7a37c2",
"assets/assets/images/CstrncpyExamples2.png": "95ce3c245260b4c537b501abd91abc59",
"assets/assets/images/playctl_next.png": "82b569d89131d85aff738a8d26d96877",
"assets/assets/images/CfgetsExample1.png": "f1d1d0aca2dcabe0ca4ee1cebbd139e6",
"assets/assets/images/CstrcpyUnboundExample3.png": "f97b8d4d0c73fe84b30835cac4d8406c",
"assets/assets/images/CstrcpyUnboundExample1.png": "0ae109f3d83db3802bd20ff1c8fcf8cb",
"assets/assets/images/CstrcatUnboundExample1.png": "0c3206228b7f06c94f5f2b81f099b1b2",
"assets/assets/images/CstrcatUnboundExample2.png": "7a30c2d182e4ee0a18199b727ee80ac0",
"assets/assets/images/splash.png": "2766030b6b31fdb9845f8006cb09bf30",
"assets/assets/images/CstrncpyExamples1.png": "c502ae1bc58c36c199afd9f671c7249d",
"assets/assets/images/pause.png": "7d19295676258e5b2944e1632164bab0",
"assets/assets/images/CstrncatExampl1.png": "1dc588c722c0f4ba7729befcee89bfc1",
"assets/assets/images/CfgetsExample2.png": "67883d0e808f114b03d4fdf53c8991ad",
"assets/assets/images/CstrcatUnboundExample3.png": "25b083b8c33049cf41399d22f281fa74",
"assets/assets/images/splash_small.png": "77d563d304b1a7b9839be0cf9398f756",
"assets/assets/images/CstrncatExample2.png": "ee74e38730035932c9bc902cb51bd920",
"assets/assets/images/CgetsExample.png": "1382405efa1cd7bf59255c6f5956b8cc",
"assets/assets/input/heap.vaccs.json": "1abbe70799edb1b8911e1201730528c2",
"assets/assets/input/t7.vaccs.json": "027f5f3e5880545143e5d2c86ae0791d",
"assets/assets/input/input.json": "c68b778f9b05869fb267c96444cc1349",
"assets/assets/html/HelpSummary/close.html": "438adcc99ec6ca56e147daf8a6d74d35",
"assets/assets/html/HelpSummary/seteuid.html": "373a818f4b3cdd425a437861106f4bfa",
"assets/assets/html/HelpSummary/write.html": "3ccb773cf7b4760b0f41fa180959fb5b",
"assets/assets/html/HelpSummary/setuid.html": "373a818f4b3cdd425a437861106f4bfa",
"assets/assets/html/HelpSummary/freopen.html": "f81f9d926f8b6c5d4d325c98d35e4737",
"assets/assets/html/HelpSummary/setresuid.html": "fca44dd88fd20d45bf84ec04f355372d",
"assets/assets/html/HelpSummary/creat.html": "9f6710f70bc5e66e461fd822e8572d68",
"assets/assets/html/HelpSummary/fread.html": "cdfeeef87a945ff6a08376102d05eb0a",
"assets/assets/html/HelpSummary/setresgid.html": "492f89d3972183a4a6fc6efcefd5dd90",
"assets/assets/html/HelpSummary/strcpy.html": "8fc64c967b2cf13ee8c274fc48e22864",
"assets/assets/html/HelpSummary/memmove.html": "bdc6e0562d92e7ad404bc4a67187a0ac",
"assets/assets/html/HelpSummary/setegid.html": "742a493432ab12ff27b5a4ebedafa031",
"assets/assets/html/HelpSummary/fclose.html": "ff40495ae1d196cb70b48c04877e3a26",
"assets/assets/html/HelpSummary/strncpy.html": "d558a97e0b351495dd3e2a385a9d9c95",
"assets/assets/html/HelpSummary/scanf.html": "6f138139be001e2216381977d8b3bf88",
"assets/assets/html/HelpSummary/setgid.html": "c2527eb9b9fcece246bfdf032c125458",
"assets/assets/html/HelpSummary/memcpy.html": "66f03a6d375da48fd59171cefcef853f",
"assets/assets/html/HelpSummary/strcat.html": "96cd689aa3d057f3d87d1dbbb602917d",
"assets/assets/html/HelpSummary/calloc.html": "f02dd9a3090379050a7ec0d2cf4e7bee",
"assets/assets/html/HelpSummary/sscanf.html": "5cd04d28f3efeebac7551acc6705ad84",
"assets/assets/html/HelpSummary/strncat.html": "53fa0830815a1e0e00bb94f8a4ff68cc",
"assets/assets/html/HelpSummary/fgets.html": "2a5c37b046d2a673d1c1fd7df17f2b24",
"assets/assets/html/HelpSummary/fstat.html": "32959e9f86e8f993605dc997960b52ea",
"assets/assets/html/HelpSummary/fopen.html": "d5d679d1da751b1d5c99c4343af30410",
"assets/assets/html/HelpSummary/_rid": "e97ede7d782243868206c27a5ecfd19c",
"assets/assets/html/HelpSummary/stat.html": "14aa744dc5c8ba88b7d31ee1a1042e05",
"assets/assets/html/HelpSummary/malloc.html": "e7a416593a1421d7de460e3ea1834650",
"assets/assets/html/HelpSummary/fdopen.html": "230c1eff62751bf5361d607328424836",
"assets/assets/html/HelpSummary/free.html": "59b56eca484e8bbd61a6bcb442001d21",
"assets/assets/html/HelpSummary/fwrite.html": "35c2fa2e018a966150a9f3e9124c60f3",
"assets/assets/html/HelpSummary/realloc.html": "52779b1416afaff98307a6a3bd63ee7e",
"assets/assets/html/HelpSummary/gets.html": "2baaae4851cee8452340e7b835eee544",
"assets/assets/html/HelpSummary/read.html": "a702b554697bcf503f30a5608d0ff7f7",
"assets/assets/html/HelpSummary/sprintf.html": "5cd04d28f3efeebac7551acc6705ad84",
"assets/assets/html/HelpSummary/printf.html": "5cd04d28f3efeebac7551acc6705ad84",
"assets/assets/html/HelpSummary/lstat.html": "6638025b209d7769002dab8e5b57c311",
"assets/assets/html/HelpSummary/open.html": "a0655d9e63a3ece38f5959b532e0f19a",
"assets/assets/html/Help/close.html": "438adcc99ec6ca56e147daf8a6d74d35",
"assets/assets/html/Help/seteuid.html": "373a818f4b3cdd425a437861106f4bfa",
"assets/assets/html/Help/write.html": "3ccb773cf7b4760b0f41fa180959fb5b",
"assets/assets/html/Help/setuid.html": "373a818f4b3cdd425a437861106f4bfa",
"assets/assets/html/Help/freopen.html": "f81f9d926f8b6c5d4d325c98d35e4737",
"assets/assets/html/Help/setresuid.html": "fca44dd88fd20d45bf84ec04f355372d",
"assets/assets/html/Help/strings.html": "50ee79021a139b38dbb09fd12a304b21",
"assets/assets/html/Help/creat.html": "9f6710f70bc5e66e461fd822e8572d68",
"assets/assets/html/Help/fread.html": "cdfeeef87a945ff6a08376102d05eb0a",
"assets/assets/html/Help/setresgid.html": "492f89d3972183a4a6fc6efcefd5dd90",
"assets/assets/html/Help/strcpy.html": "e683f9a00f42b19b021a7b1aba30b073",
"assets/assets/html/Help/memmove.html": "77d44cd16ff374fd7731824bdba2041e",
"assets/assets/html/Help/setegid.html": "742a493432ab12ff27b5a4ebedafa031",
"assets/assets/html/Help/fclose.html": "ff40495ae1d196cb70b48c04877e3a26",
"assets/assets/html/Help/strncpy.html": "92733382c2f3cc89c0614478cd6c3131",
"assets/assets/html/Help/setgid.html": "c2527eb9b9fcece246bfdf032c125458",
"assets/assets/html/Help/memcpy.html": "6ef30c7af69642eb183a281e90e26f61",
"assets/assets/html/Help/strcat.html": "3a6bd5df43af1a686742b1381e6dc24d",
"assets/assets/html/Help/calloc.html": "85fdc515ae2f4129c71e8a5e8e6d9380",
"assets/assets/html/Help/strncat.html": "4c0cf73c11c503ff0b5691a61542a03e",
"assets/assets/html/Help/fgets.html": "47df63de60e7ec93ac3df1ef15ae2fc3",
"assets/assets/html/Help/fstat.html": "32959e9f86e8f993605dc997960b52ea",
"assets/assets/html/Help/fopen.html": "d5d679d1da751b1d5c99c4343af30410",
"assets/assets/html/Help/_rid": "8ecc1d8829b6f17da477b57774ac0210",
"assets/assets/html/Help/stat.html": "14aa744dc5c8ba88b7d31ee1a1042e05",
"assets/assets/html/Help/malloc.html": "5e013f9d04e2292f4b67cd414b738e07",
"assets/assets/html/Help/fdopen.html": "230c1eff62751bf5361d607328424836",
"assets/assets/html/Help/images/jpg/CstrcpyUnboundExample2.jpeg": "4cfeaa4bf59c330ba3d53e6c7c0de716",
"assets/assets/html/Help/images/jpg/CstrncpyExamples2.jpeg": "217c7721f68d97eaffcf21e97460043f",
"assets/assets/html/Help/images/jpg/CfgetsExample2.jpeg": "f9dae635e0708f698e9988bcba70ae5a",
"assets/assets/html/Help/images/jpg/CstrcatUnboundExample1.jpeg": "10a7bee6f351b3b0e05cc5f9a963ebf6",
"assets/assets/html/Help/images/jpg/CgetsExample.jpeg": "099ef6a33d86b4b7711f17891a4568f1",
"assets/assets/html/Help/images/jpg/CmemExample.jpeg": "75305aed63ef07d83a6e223541742902",
"assets/assets/html/Help/images/jpg/CStringFigureSecurity.jpeg": "644c8fe81a8181082b51f782937be5c3",
"assets/assets/html/Help/images/jpg/CstrncpyExamples1.jpeg": "7abb51f0fbd684225a39bcfa344322ee",
"assets/assets/html/Help/images/jpg/_rid": "03438082d7e7f67343ed3f21ea341dce",
"assets/assets/html/Help/images/jpg/CstrcpyUnboundExample3.jpeg": "722d92d3eeea6e9371c0b028b0536552",
"assets/assets/html/Help/images/jpg/CstrcatUnboundExample3.jpeg": "019b95350315f6dbc9a3d3565505c629",
"assets/assets/html/Help/images/jpg/CstrncatExampl1.jpeg": "ab55ea4099fd2285a519365545ecba26",
"assets/assets/html/Help/images/jpg/CstrncatExample2.jpeg": "64b85d18b158f5ad8f7e88bf6192bfd2",
"assets/assets/html/Help/images/jpg/CstrcatUnboundExample2.jpeg": "952ceafd9fd7aacb16bd3c13f5381d0d",
"assets/assets/html/Help/images/jpg/CstrcpyUnboundExample1.jpg": "caf0cb6a495a8be08f53daca880a998c",
"assets/assets/html/Help/images/jpg/CfgetsExample1.jpeg": "82221143dff0fda22e42321f53fe205c",
"assets/assets/html/Help/free.html": "59b56eca484e8bbd61a6bcb442001d21",
"assets/assets/html/Help/fwrite.html": "35c2fa2e018a966150a9f3e9124c60f3",
"assets/assets/html/Help/realloc.html": "52779b1416afaff98307a6a3bd63ee7e",
"assets/assets/html/Help/gets.html": "4ec9dbe80ac7fd70b80a1395da99d091",
"assets/assets/html/Help/read.html": "a702b554697bcf503f30a5608d0ff7f7",
"assets/assets/html/Help/lstat.html": "6638025b209d7769002dab8e5b57c311",
"assets/assets/html/Help/open.html": "a0655d9e63a3ece38f5959b532e0f19a",
"assets/AssetManifest.bin": "c75646b638f376367f428019a4491e96",
"assets/packages/cupertino_icons/assets/CupertinoIcons.ttf": "33b7d9392238c04c131b6ce224e13711",
"assets/NOTICES": "02a415f4ed609af941b40f380c1c5739",
"assets/AssetManifest.json": "3d0cc1c81c964ec95e1b3596cc80c511",
"assets/AssetManifest.bin.json": "e6285f30ad83b4d30ee05f7d31161e70",
"assets/shaders/ink_sparkle.frag": "ecc85a2e95f5e9f53123dcaf8cb9b6ce",
"assets/fonts/MaterialIcons-Regular.otf": "1d82e14f9fa711be560627294de45719",
"assets/FontManifest.json": "dc3d03800ccca4601324923c0b1d6d57",
"main.dart.js": "ee88894df90e4ada9897a20998bf726d",
"icons/Icon-maskable-192.png": "c457ef57daa1d16f64b27b786ec2ea3c",
"icons/Icon-maskable-512.png": "301a7604d45b3e739efc881eb04896ea",
"icons/Icon-512.png": "96e752610906ba2a93c65f8abe1645f1",
"icons/Icon-192.png": "ac9a721a12bbc803b44f645561ecb1e1",
"manifest.json": "32e055ef95571e9680c04ab3dfc35ecb",
"canvaskit/skwasm.js": "1ef3ea3a0fec4569e5d531da25f34095",
"canvaskit/canvaskit.js": "140ccb7d34d0a55065fbd422b843add6",
"canvaskit/canvaskit.js.symbols": "58832fbed59e00d2190aa295c4d70360",
"canvaskit/skwasm_heavy.js.symbols": "3c01ec03b5de6d62c34e17014d1decd3",
"canvaskit/skwasm.js.symbols": "0088242d10d7e7d6d2649d1fe1bda7c1",
"canvaskit/skwasm.wasm": "264db41426307cfc7fa44b95a7772109",
"canvaskit/canvaskit.wasm": "07b9f5853202304d3b0749d9306573cc",
"canvaskit/skwasm_heavy.js": "413f5b2b2d9345f37de148e2544f584f",
"canvaskit/skwasm_heavy.wasm": "8034ad26ba2485dab2fd49bdd786837b",
"canvaskit/chromium/canvaskit.js": "5e27aae346eee469027c80af0751d53d",
"canvaskit/chromium/canvaskit.js.symbols": "193deaca1a1424049326d4a91ad1d88d",
"canvaskit/chromium/canvaskit.wasm": "24c77e750a7fa6d474198905249ff506",
"version.json": "6aa3ff86f7c94feb6ee75ef653e94347",
"flutter_bootstrap.js": "88ad2fa4152d2068147e1e6310432c48",
"index.html": "adb22d753d791c149782089281a0ece2",
"/": "adb22d753d791c149782089281a0ece2"};
// The application shell files that are downloaded before a service worker can
// start.
const CORE = ["main.dart.js",
"index.html",
"flutter_bootstrap.js",
"assets/AssetManifest.bin.json",
"assets/FontManifest.json"];

// During install, the TEMP cache is populated with the application shell files.
self.addEventListener("install", (event) => {
  self.skipWaiting();
  return event.waitUntil(
    caches.open(TEMP).then((cache) => {
      return cache.addAll(
        CORE.map((value) => new Request(value, {'cache': 'reload'})));
    })
  );
});
// During activate, the cache is populated with the temp files downloaded in
// install. If this service worker is upgrading from one with a saved
// MANIFEST, then use this to retain unchanged resource files.
self.addEventListener("activate", function(event) {
  return event.waitUntil(async function() {
    try {
      var contentCache = await caches.open(CACHE_NAME);
      var tempCache = await caches.open(TEMP);
      var manifestCache = await caches.open(MANIFEST);
      var manifest = await manifestCache.match('manifest');
      // When there is no prior manifest, clear the entire cache.
      if (!manifest) {
        await caches.delete(CACHE_NAME);
        contentCache = await caches.open(CACHE_NAME);
        for (var request of await tempCache.keys()) {
          var response = await tempCache.match(request);
          await contentCache.put(request, response);
        }
        await caches.delete(TEMP);
        // Save the manifest to make future upgrades efficient.
        await manifestCache.put('manifest', new Response(JSON.stringify(RESOURCES)));
        // Claim client to enable caching on first launch
        self.clients.claim();
        return;
      }
      var oldManifest = await manifest.json();
      var origin = self.location.origin;
      for (var request of await contentCache.keys()) {
        var key = request.url.substring(origin.length + 1);
        if (key == "") {
          key = "/";
        }
        // If a resource from the old manifest is not in the new cache, or if
        // the MD5 sum has changed, delete it. Otherwise the resource is left
        // in the cache and can be reused by the new service worker.
        if (!RESOURCES[key] || RESOURCES[key] != oldManifest[key]) {
          await contentCache.delete(request);
        }
      }
      // Populate the cache with the app shell TEMP files, potentially overwriting
      // cache files preserved above.
      for (var request of await tempCache.keys()) {
        var response = await tempCache.match(request);
        await contentCache.put(request, response);
      }
      await caches.delete(TEMP);
      // Save the manifest to make future upgrades efficient.
      await manifestCache.put('manifest', new Response(JSON.stringify(RESOURCES)));
      // Claim client to enable caching on first launch
      self.clients.claim();
      return;
    } catch (err) {
      // On an unhandled exception the state of the cache cannot be guaranteed.
      console.error('Failed to upgrade service worker: ' + err);
      await caches.delete(CACHE_NAME);
      await caches.delete(TEMP);
      await caches.delete(MANIFEST);
    }
  }());
});
// The fetch handler redirects requests for RESOURCE files to the service
// worker cache.
self.addEventListener("fetch", (event) => {
  if (event.request.method !== 'GET') {
    return;
  }
  var origin = self.location.origin;
  var key = event.request.url.substring(origin.length + 1);
  // Redirect URLs to the index.html
  if (key.indexOf('?v=') != -1) {
    key = key.split('?v=')[0];
  }
  if (event.request.url == origin || event.request.url.startsWith(origin + '/#') || key == '') {
    key = '/';
  }
  // If the URL is not the RESOURCE list then return to signal that the
  // browser should take over.
  if (!RESOURCES[key]) {
    return;
  }
  // If the URL is the index.html, perform an online-first request.
  if (key == '/') {
    return onlineFirst(event);
  }
  event.respondWith(caches.open(CACHE_NAME)
    .then((cache) =>  {
      return cache.match(event.request).then((response) => {
        // Either respond with the cached resource, or perform a fetch and
        // lazily populate the cache only if the resource was successfully fetched.
        return response || fetch(event.request).then((response) => {
          if (response && Boolean(response.ok)) {
            cache.put(event.request, response.clone());
          }
          return response;
        });
      })
    })
  );
});
self.addEventListener('message', (event) => {
  // SkipWaiting can be used to immediately activate a waiting service worker.
  // This will also require a page refresh triggered by the main worker.
  if (event.data === 'skipWaiting') {
    self.skipWaiting();
    return;
  }
  if (event.data === 'downloadOffline') {
    downloadOffline();
    return;
  }
});
// Download offline will check the RESOURCES for all files not in the cache
// and populate them.
async function downloadOffline() {
  var resources = [];
  var contentCache = await caches.open(CACHE_NAME);
  var currentContent = {};
  for (var request of await contentCache.keys()) {
    var key = request.url.substring(origin.length + 1);
    if (key == "") {
      key = "/";
    }
    currentContent[key] = true;
  }
  for (var resourceKey of Object.keys(RESOURCES)) {
    if (!currentContent[resourceKey]) {
      resources.push(resourceKey);
    }
  }
  return contentCache.addAll(resources);
}
// Attempt to download the resource online before falling back to
// the offline cache.
function onlineFirst(event) {
  return event.respondWith(
    fetch(event.request).then((response) => {
      return caches.open(CACHE_NAME).then((cache) => {
        cache.put(event.request, response.clone());
        return response;
      });
    }).catch((error) => {
      return caches.open(CACHE_NAME).then((cache) => {
        return cache.match(event.request).then((response) => {
          if (response != null) {
            return response;
          }
          throw error;
        });
      });
    })
  );
}
