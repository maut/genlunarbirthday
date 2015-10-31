 <td id="wikicontent" class="psdescription">
 <h2><a name=""></a>GenLunarBirthday</h2><p>Generate birthday dates base on lunar birthdays for google calendar import Can be used for notifying birthdays using google calendar </p><p>根据农历生日生成可用于谷歌日历导入的csv文件, 然后导入谷歌日历就可以每年收到所有家人的生日提醒了 </p><p>Generate birthdays up to 2100 </p><p>生成至2100年的生日 </p><ul><li><a href="http://code.google.com/p/genlunarbirthday/downloads/list" rel="nofollow">下载</a> </li><li>从源文件编译: <a href="http://golang.org/doc/install#install" rel="nofollow">下载安装Go</a>; <tt>go get code.google.com/p/genlunarbirthday</tt> </li></ul><p>用法(Windows): <ol><li>编辑days.txt文件如下: (-4 代表闰四月) </li></p><pre class="prettyprint">  父亲生日: 1959.01.01
  母亲生日: 1963.-4.01</pre><li>将days.txt拖到genlunarbirthday.exe上 </li><li>或者在命令行输入genlunarbirthday.exe days.txt </li><li>应该会有一个import.csv生成 </li><li>可以选择编辑这个文件, 但是要用文本编辑器打开 </li><li>在谷歌日历里导入并设置提醒 </li><li>每年收到邮件/短信/等等生日提醒 </li></ol><p>用法2: <ol><li>打开 <a href="http://play.golang.org/p/pjHLGH_HjP" rel="nofollow">http://play.golang.org/p/pjHLGH_HjP</a> </li><li>修改<tt>const example</tt>的内容 </li><li>点击<tt>Run</tt> </li><li>复制粘贴结果到Excel或文本编辑器 </li></p><pre class="prettyprint">Usage:
    genlunarbirthday.exe birthdays.txt
Then a import.csv will be created
Example birthday.txt:
&quot;father&#x27;s birthday&quot;: 1959.01.01
&quot;mother&#x27;s birthday&quot;: 1963.-4.01

Subject, Start Date, Start Time, End Date, End Time, Private, All Day Event, Location
&quot;father&#x27;s birthday&quot;, 1959-2-8, 8:00 AM, 1959-2-8, , FALSE, TRUE, HOME
&quot;father&#x27;s birthday&quot;, 1960-1-28, 8:00 AM, 1960-1-28, , FALSE, TRUE, HOME
&quot;father&#x27;s birthday&quot;, 1961-2-15, 8:00 AM, 1961-2-15, , FALSE, TRUE, HOME
...
&quot;father&#x27;s birthday&quot;, 2095-2-5, 8:00 AM, 2095-2-5, , FALSE, TRUE, HOME
&quot;father&#x27;s birthday&quot;, 2096-1-25, 8:00 AM, 2096-1-25, , FALSE, TRUE, HOME
&quot;father&#x27;s birthday&quot;, 2097-2-12, 8:00 AM, 2097-2-12, , FALSE, TRUE, HOME
&quot;father&#x27;s birthday&quot;, 2098-2-1, 8:00 AM, 2098-2-1, , FALSE, TRUE, HOME
&quot;father&#x27;s birthday&quot;, 2099-1-21, 8:00 AM, 2099-1-21, , FALSE, TRUE, HOME
&quot;father&#x27;s birthday&quot;, 2100-2-9, 8:00 AM, 2100-2-9, , FALSE, TRUE, HOME
...
&quot;mother&#x27;s birthday&quot;, 1963-5-23, 8:00 AM, 1963-5-23, , FALSE, TRUE, HOME
&quot;mother&#x27;s birthday&quot;, 1974-5-22, 8:00 AM, 1974-5-22, , FALSE, TRUE, HOME
...
&quot;mother&#x27;s birthday&quot;, 2088-5-21, 8:00 AM, 2088-5-21, , FALSE, TRUE, HOME
&quot;mother&#x27;s birthday&quot;, 2096-5-22, 8:00 AM, 2096-5-22, , FALSE, TRUE, HOME</pre></ol>
 </td>
