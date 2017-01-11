XCSS
====

# 概要
　UXFrameworkのlayout.xmlファイルを書くにあたって、HTMLに対するCSSのような記法を使用することができるシステムです。記述の重複を減らし、変更を容易にします。

　CSSの知識は必須ではありませんが。id, classによるセレクタの設定方法などを参考にしているため、知っていた方が理解が早まります。

# 特別なファイル
　XCSSは二種類のファイルを特別視して扱います。

## XCSSファイル
　"xcss.xml"というファイル名のXMLファイルを、XCSS(XML's css)ファイルとして扱います。XCSSファイルは、HTMLにおけるCSSファイルに相当するものです。変換に必要な設定を記述します。SXMLファイルをXMLファイルに変換する際に適用されます。

　あるXCSSファイルの設定は、そのファイルが存在するフォルダと、それ以下のフォルダにあるSXMLファイルに対して適用されます。ひとつのフォルダに複数のXCSSファイルが存在する場合、ファイル名順に適用されます。

### XCSSファイルの記述法
　XCSSファイルは形式的にはxmlファイルです。ルートは必ずstylesタグとし、その直下に存在するタグひとつが、ひとつのスタイル指定となります。

　XCSSはタグ名および"type", "id", "class"の属性

    minimum_xcss.xml
    <styles>
      <item type="image" id="background" w="800" h="480" imageurl="@DNA/common/bg.png"/>
    </styles>

## SXMLファイル
　"sxml.xml"というファイル名のXMLファイルを、SXML(Styled XML)ファイルとして扱います。SXMLファイルは、XCSSファイルからの設定を適用された後、".xml"ファイルという名前のXMLファイルに変換されます。

     minimum_sxml.xml
     <pane>
       <widgets>
         <item type="image" id="background" />
       </widgets>
     </pane>

     minimum.xml
     <pane>
       <widgets>
         <item type="image" id="background" w="800" h="480" imageurl="@DNA/common/bg.png"></item>
       </widgets>
     </pane>


# フラグ
      -c string
            short form of "class"
      -class string
            classes apply to all elements. separator is space e.g. "foo bar"
      -d    short form of "debug"
      -debug
            add comment where attribute's value came from
      -r string
            short form of "root" (default ".")
      -root string
            path to start recursive walk (default ".")
      -w    short form of "watch"
      -watch
            watch "xcss" and "sxml" files and run convert when these files are changed


## class(c) string
　フラグとしてclassを設定します。ここで設定されたclassは全ての要素に適用されます。ビルド時に設定を分ける場合に使用することを想定しています。複数設定する場合の区切り文字はスペースです。たとえばfooとbarふたつのclassを設定するときは"foo bar"となります。

## debug(d)
　debugフラグを設定すると、出力されたXMLにデバッグ用のコメントを出力します。タグ・アトリビュートが元のSXMLファイルに直接指定されたものではない場合、どのXCSSファイルのどのセレクタから来たものかがわかります。

## root(r) string
　変換対象となるルートディレクトリを指定します。ここ以下のパスが対象になります。デフォルトは"."です。つまり実行ファイルのカレントディレクトリとなります。

## watch(w)
　watchフラグを指定すると一度全体を変換した後、終了せずにファイルの変更を監視します。SXMLファイルの変更があった場合、そのファイルだけを再度変換してXMLファイルを出力します。XCSSファイルの追加・削除・変更、SXMLファイルの追加・削除があった場合は、再度全体を変換します。


# インストール
　環境に合わせたバイナリを以下から取得し配置して下さい。

# 作者
Takanori Kido <tkido@uievolution.com>
