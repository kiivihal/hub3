// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: https://codemirror.net/LICENSE
(function(mod){if("object"==typeof exports&&"object"==typeof module)// CommonJS
mod(require("../../lib/codemirror"));else if("function"==typeof define&&define.amd)// AMD
define(["../../lib/codemirror"],mod);else// Plain browser env
mod(CodeMirror)})(function(CodeMirror){"use strict";CodeMirror.defineMode("groovy",function(config){function words(str){for(var obj={},words=str.split(" "),i=0;i<words.length;++i)obj[words[i]]=!0/* ignoreName */ /* skipSlots */;return obj}var keywords=words("abstract as assert boolean break byte case catch char class const continue def default "+"do double else enum extends final finally float for goto if implements import in "+"instanceof int interface long native new package private protected public return "+"short static strictfp super switch synchronized threadsafe throw throws trait transient "+"try void volatile while"),blockKeywords=words("catch class def do else enum finally for if interface switch trait try while"),standaloneKeywords=words("return break continue"),atoms=words("null true false this"),curPunc;function tokenBase(stream,state){var ch=stream.next();if("\""==ch||"'"==ch){return startString(ch,stream,state)}if(/[\[\]{}\(\),;\:\.]/.test(ch)){curPunc=ch;return null}if(/\d/.test(ch)){stream.eatWhile(/[\w\.]/);if(stream.eat(/eE/)){stream.eat(/\+\-/);stream.eatWhile(/\d/)}return"number"}if("/"==ch){if(stream.eat("*")){state.tokenize.push(tokenComment);return tokenComment(stream,state)}if(stream.eat("/")){stream.skipToEnd();return"comment"}if(expectExpression(state.lastToken,/* ignoreName */ /* eat */!1/* skipSlots */ /* skipSlots */)){return startString(ch,stream,state)}}if("-"==ch&&stream.eat(">")){curPunc="->";return null}if(/[+\-*&%=<>!?|\/~]/.test(ch)){stream.eatWhile(/[+\-*&%=<>|~]/);return"operator"}stream.eatWhile(/[\w\$_]/);if("@"==ch){stream.eatWhile(/[\w\$_\.]/);return"meta"}if("."==state.lastToken)return"property";if(stream.eat(":")){curPunc="proplabel";return"property"}var cur=stream.current();if(atoms.propertyIsEnumerable(cur)){return"atom"}if(keywords.propertyIsEnumerable(cur)){if(blockKeywords.propertyIsEnumerable(cur))curPunc="newstatement";else if(standaloneKeywords.propertyIsEnumerable(cur))curPunc="standalone";return"keyword"}return"variable"}tokenBase.isBase=!0;function startString(quote,stream,state){var tripleQuoted=!1;if("/"!=quote&&stream.eat(quote)){if(stream.eat(quote))tripleQuoted=!0;else return"string"}function t(stream,state){var escaped=!1,next,end=!tripleQuoted;while(null!=(next=stream.next())){if(next==quote&&!escaped){if(!tripleQuoted){break}if(stream.match(quote+quote)){end=!0;break}}if("\""==quote&&"$"==next&&!escaped&&stream.eat("{")){state.tokenize.push(tokenBaseUntilBrace());return"string"}escaped=!escaped&&"\\"==next}if(end)state.tokenize.pop();return"string"}state.tokenize.push(t);return t(stream,state)}function tokenBaseUntilBrace(){var depth=1;function t(stream,state){if("}"==stream.peek()){depth--;if(0==depth){state.tokenize.pop();return state.tokenize[state.tokenize.length-1](stream,state)}}else if("{"==stream.peek()){depth++}return tokenBase(stream,state)}t.isBase=!0;return t}function tokenComment(stream,state){var maybeEnd=!1,ch;while(ch=stream.next()){if("/"==ch&&maybeEnd){state.tokenize.pop();break}maybeEnd="*"==ch}return"comment"}function expectExpression(last,newline){return!last||"operator"==last||"->"==last||/[\.\[\{\(,;:]/.test(last)||"newstatement"==last||"keyword"==last||"proplabel"==last||"standalone"==last&&!newline}function Context(indented,column,type,align,prev){this.indented=indented;this.column=column;this.type=type;this.align=align;this.prev=prev}function pushContext(state,col,type){return state.context=new Context(state.indented,col,type,null,state.context)}function popContext(state){var t=state.context.type;if(")"==t||"]"==t||"}"==t)state.indented=state.context.indented;return state.context=state.context.prev}// Interface
return{startState:function(basecolumn){return{tokenize:[tokenBase],context:new Context((basecolumn||0)-config.indentUnit,0,"top",!1),indented:0,startOfLine:!0,lastToken:null}},token:function(stream,state){var ctx=state.context;if(stream.sol()){if(null==ctx.align)ctx.align=!1;state.indented=stream.indentation();state.startOfLine=!0;// Automatic semicolon insertion
if("statement"==ctx.type&&!expectExpression(state.lastToken,!0)){popContext(state);ctx=state.context}}if(stream.eatSpace())return null;curPunc=null;var style=state.tokenize[state.tokenize.length-1](stream,state);if("comment"==style)return style;if(null==ctx.align)ctx.align=!0;if((";"==curPunc||":"==curPunc)&&"statement"==ctx.type)popContext(state);// Handle indentation for {x -> \n ... }
else if("->"==curPunc&&"statement"==ctx.type&&"}"==ctx.prev.type){popContext(state);state.context.align=!1}else if("{"==curPunc)pushContext(state,stream.column(),"}");else if("["==curPunc)pushContext(state,stream.column(),"]");else if("("==curPunc)pushContext(state,stream.column(),")");else if("}"==curPunc){while("statement"==ctx.type)ctx=popContext(state);if("}"==ctx.type)ctx=popContext(state);while("statement"==ctx.type)ctx=popContext(state)}else if(curPunc==ctx.type)popContext(state);else if("}"==ctx.type||"top"==ctx.type||"statement"==ctx.type&&"newstatement"==curPunc)pushContext(state,stream.column(),"statement");state.startOfLine=!1;state.lastToken=curPunc||style;return style},indent:function(state,textAfter){if(!state.tokenize[state.tokenize.length-1].isBase)return CodeMirror.Pass;var firstChar=textAfter&&textAfter.charAt(0),ctx=state.context;if("statement"==ctx.type&&!expectExpression(state.lastToken,!0))ctx=ctx.prev;var closing=firstChar==ctx.type;if("statement"==ctx.type)return ctx.indented+("{"==firstChar?0:config.indentUnit);else if(ctx.align)return ctx.column+(closing?0:1);else return ctx.indented+(closing?0:config.indentUnit)},electricChars:"{}",closeBrackets:{triples:"'\""},fold:"brace",blockCommentStart:"/*",blockCommentEnd:"*/",lineComment:"//"}});CodeMirror.defineMIME("text/x-groovy","groovy")});