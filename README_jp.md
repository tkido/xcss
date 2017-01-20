XCSS
====

# 概要
　UXFrameworkのlayout.xmlファイルを書くにあたって、HTMLに対するCSSのような記法を使用することができるシステムです。記述の重複を減らし、管理・変更を容易にします。

　CSSの知識は必須ではありませんが。id, classによるセレクタの設定方法などを参考にしているため、知っていた方が理解が早まります。

# 特別なファイル
　XCSSは二種類のxmlファイルを特別視して扱います。

## XCSSファイル
　ファイル名の末尾が"\*\_xcss.xml"であるXMLファイルを、XCSS(XML's css)ファイルとして扱います。XCSSファイルは、HTMLにおけるCSSファイルに相当するもので、変換に必要な設定を記述します。

　あるXCSSファイルの設定は、そのファイルが存在するフォルダと、それ以下のフォルダにあるSXMLファイルに対して適用されます。ひとつのフォルダに複数のXCSSファイルが存在する場合、ファイル名順に適用されます。

### XCSSファイルの記述法
　XCSSファイルは形式的にはxmlファイルです。ルートは必ずstylesタグとし、その直下に存在するタグひとつが、ひとつのスタイル指定となります。

　XCSSはタグ名および"type", "id", "class"の属性を、設定を適用するタグを決定するセレクタとして利用します。

~~~
<!-- minimum_xcss.xml -->
<styles>
  <item type="image" id="background" w="800" h="480" imageurl="@DNA/common/bg.png"/>
</styles>
~~~

## SXMLファイル
　ファイル名の末尾が"\*\_sxml.xml"であるXMLファイルを、SXML(Styled XML)ファイルとして扱います。SXMLファイルは、XCSSファイルからの設定を適用された後、"\*.xml"という名前のXMLファイルに変換されます。

変換前：
~~~
<!-- minimum_sxml.xml -->
 <pane>
   <widgets>
     <item type="image" id="background" />
   </widgets>
 </pane>
 ~~~
変換後：
 ~~~
 <!-- minimum.xml -->
 <pane>
   <widgets>
     <item type="image" id="background" w="800" h="480" imageurl="@DNA/common/bg.png"></item>
   </widgets>
 </pane>
 ~~~


# フラグ
~~~
-c string
	short form of "class"
-class string
	classes apply to all elements. separator is space e.g. "foo bar"
-d
  short form of "debug"
-debug
	add comment where attribute's value came from
-delete
	delete all XCSS and SXML files after conversion
-r string
	short form of "root" (default ".")
-root string
	path to start recursive walk (default ".")
-w
  short form of "watch"
-watch
	watch "xcss" and "sxml" files and run convert when these files are changed
~~~

## class(c) string
　フラグとしてclassを設定します。ここで設定されたclassは、全ての要素に適用されます。ビルド時に設定を分けるような場合に使用します。複数設定する場合の区切り文字はスペースです。たとえばfooとbarふたつのclassを設定するときは"foo bar"となります。

## debug(d)
　debugフラグを設定すると、出力されたXMLにデバッグ用のコメントを出力します。タグ・アトリビュートが元のSXMLファイルに直接指定されたものではない場合、どのXCSSファイルのどのセレクタから来たものかがわかります。

## delete
　deleteフラグを指定すると、変換を適用した後、全てのXCSSおよびSXMLファイルを削除します。リリースバイナリのビルド時など、変換後のXMLのみを残したい場合に利用します。

## root(r) string
　変換対象となるルートディレクトリを指定します。ここ以下のパスのみが対象になります。デフォルトは"."、つまり実行ファイルのカレントディレクトリとなります。

## watch(w)
　watchフラグを指定すると一度全体を変換した後、プログラムを終了せずにファイルの変更を監視します。SXMLファイルの変更があった場合、そのファイルだけを再度変換してXMLファイルを出力します。XCSSファイルの追加・削除・変更、SXMLファイルの追加・削除があった場合は、再度全体を変換します。


# インストール
　環境に合わせたバイナリを以下から取得し配置して下さい。

# 作者
Takanori Kido <tkido@uievolution.com>
