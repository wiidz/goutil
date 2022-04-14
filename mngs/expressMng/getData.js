// ==UserScript==
// @name         New Userscript
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  try to take over the world!
// @author       You
// @match        https://market.aliyun.com/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=aliyun.com
// @grant        none
// ==/UserScript==

(function() {
    'use strict';
    console.log("start");
    // Your code here...
    setTimeout(function(){


        let expressData = [];
        let trs = document.querySelector(".Table").children[0].children;
        console.log(trs);
        for(let i=1;i<trs.length;i++){
            expressData.push({"code":dropSymbol(trs[i].children[1].textContent),"name":dropSymbol(trs[i].children[0].textContent)});
            expressData.push({"code":dropSymbol(trs[i].children[3].textContent),"name":dropSymbol(trs[i].children[2].textContent)});
        }
        console.log("expressData",expressData);
        getJson(expressData);
        getGoConst(expressData);
    },1000);

    function dropSymbol(str){
        return str.replaceAll("\n","").replaceAll("\t","");
    }

    function getJson(data){
        console.log(JSON.stringify(data));
    }

    function getGoConst(data){
        let temp = "type ExpressKind string\n";
        for(let i = 0;i<data.length;i++){
            temp += "const " + data[i].code + " ExpressKind  = \"" + data[i].name + "\"\n";
        }
        console.log(temp);
    }

})();


