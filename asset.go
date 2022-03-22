package jd_cookie

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/buger/jsonparser"
	"github.com/cdle/sillyGirl/core"
	"github.com/cdle/sillyGirl/develop/qinglong"
)

type JdCookie struct {
	ID        int
	PtKey     string
	PtPin     string
	WsKey     string
	Note      string
	Nickname  string
	BeanNum   string
	UserLevel string
	LevelName string
}

//增加随机ua
var USER_AGENTS = []string{
	`jdltapp;iPad;3.7.0;14.4;network/wifi;Mozilla/5.0 (iPad; CPU OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;10;2346663656561603-4353564623932316;network/wifi;model/ONEPLUS A5010;addressid/0;aid/2dfceea045ed292a;oaid/;osVer/29;appBuild/1436;psn/BS6Y9SAiw0IpJ4ro7rjSOkCRZTgR3z2K|10;psq/5;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/10.5;jdv/;ref/com.jd.jdlite.lib.personal.view.fragment.JDPersonalFragment;partner/oppo;apprpd/MyJD_Main;eufv/1;Mozilla/5.0 (Linux; Android 10; ONEPLUS A5010 Build/QKQ1.191014.012; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045140 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.1;59d6ae6e8387bd09fe046d5b8918ead51614e80a;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone12,1;hasOCPay/0;appBuild/1017;supportBestPay/0;addressid/;pv/1.26;apprpd/;ref/JDLTSubMainPageViewController;psq/0;ads/;psn/59d6ae6e8387bd09fe046d5b8918ead51614e80a|3;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.1;Mozilla/5.0 (iPhone; CPU iPhone OS 14_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;13.5;22d679c006bf9c087abf362cf1d2e0020ebb8798;network/wifi;ADID/10857A57-DDF8-4A0D-A548-7B8F43AC77EE;hasUPPay/0;pushNoticeIsOpen/1;lang/zh_CN;model/iPhone12,1;addressid/2378947694;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/15.7;apprpd/Allowance_Registered;ref/JDLTTaskCenterViewController;psq/6;ads/;psn/22d679c006bf9c087abf362cf1d2e0020ebb8798|22;jdv/0|kong|t_1000170135|tuiguang|notset|1614153044558|1614153044;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 13.5;Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;10;2616935633265383-5333463636261326;network/UNKNOWN;model/M2007J3SC;addressid/1840745247;aid/ba9e3b5853dccb1b;oaid/371d8af7dd71e8d5;osVer/29;appBuild/1436;psn/t7JmxZUXGkimd4f9Jdul2jEeuYLwxPrm|8;psq/6;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/5.6;jdv/;ref/com.jd.jdlite.lib.jdlitemessage.view.activity.MessageCenterMainActivity;partner/xiaomi;apprpd/MessageCenter_MessageMerge;eufv/1;Mozilla/5.0 (Linux; Android 10; M2007J3SC Build/QKQ1.200419.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045135 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.3;d7beab54ae7758fa896c193b49470204fbb8fce9;network/4g;ADID/97AD46C9-6D49-4642-BF6F-689256673906;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/6.28;apprpd/;ref/JDLTRedPacketViewController;psq/3;ads/;psn/d7beab54ae7758fa896c193b49470204fbb8fce9|8;jdv/0|kong|t_1001707023_|jingfen|79ad0319fa4d47e38521a616d80bc4bd|1613800945610|1613824900;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;9;D246836333735-3264353430393;network/4g;model/MIX 2;addressid/138678023;aid/bf8bcf1214b3832a;oaid/308540d1f1feb2f5;osVer/28;appBuild/1436;psn/Z/rGqfWBY/h5gcGFnVIsRw==|16;psq/3;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 9;osv/9;pv/13.7;jdv/;ref/com.jd.jdlite.lib.personal.view.fragment.JDPersonalFragment;partner/xiaomi;apprpd/MyJD_Main;eufv/1;Mozilla/5.0 (Linux; Android 9; MIX 2 Build/PKQ1.190118.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045135 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.4;eb5a9e7e596e262b4ffb3b6b5c830984c8a5c0d5;network/wifi;ADID/5603541B-30C1-4B5C-A782-20D0B569D810;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone9,2;addressid/1041002757;hasOCPay/0;appBuild/101;supportBestPay/0;pv/34.6;apprpd/MyJD_Main;ref/MyJdMTAManager;psq/5;ads/;psn/eb5a9e7e596e262b4ffb3b6b5c830984c8a5c0d5|44;jdv/0|androidapp|t_335139774|appshare|CopyURL|1612612940307|1612612944;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;21631ed983b3e854a3154b0336413825ad0d6783;network/3g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone13,4;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.47;apprpd/;ref/JDLTSubMainPageViewController;psq/8;ads/;psn/21631ed983b3e854a3154b0336413825ad0d6783|9;jdv/0|direct|-|none|-|1614150725100|1614225882;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;13.5;500a795cb2abae60b877ee4a1930557a800bef1c;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone8,1;addressid/669949466;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/9.11;apprpd/;ref/JDLTSubMainPageViewController;psq/10;ads/;psn/500a795cb2abae60b877ee4a1930557a800bef1c|11;jdv/;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 13.5;Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPad;3.7.0;14.4;f5e7b7980fb50efc9c294ac38653c1584846c3db;network/wifi;hasUPPay/0;pushNoticeIsOpen/1;lang/zh_CN;model/iPad6,3;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/231.11;pap/JA2020_3112531|3.7.0|IOS 14.4;apprpd/;psn/f5e7b7980fb50efc9c294ac38653c1584846c3db|305;usc/kong;jdv/0|kong|t_1000170135|tuiguang|notset|1613606450668|1613606450;umd/tuiguang;psq/2;ucp/t_1000170135;app_device/IOS;utr/notset;ref/JDLTRedPacketViewController;adk/;ads/;Mozilla/5.0 (iPad; CPU OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;19fef5419f88076c43f5317eabe20121d52c6a61;network/wifi;ADID/00000000-0000-0000-0000-000000000000;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,8;addressid/3430850943;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/10.4;apprpd/;ref/JDLTSubMainPageViewController;psq/3;ads/;psn/19fef5419f88076c43f5317eabe20121d52c6a61|16;jdv/0|kong|t_1001327829_|jingfen|f51febe09dd64b20b06bc6ef4c1ad790#/|1614096460311|1614096511;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148`,
	`jdltapp;iPhone;3.7.0;12.2;f995bc883282f7c7ea9d7f32da3f658127aa36c7;network/4g;ADID/9F40F4CA-EA7C-4F2E-8E09-97A66901D83E;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone10,4;addressid/525064695;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/11.11;apprpd/;ref/JDLTSubMainPageViewController;psq/2;ads/;psn/f995bc883282f7c7ea9d7f32da3f658127aa36c7|22;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 12.2;Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;10;5366566313931326-6633931643233693;network/wifi;model/Mi9 Pro 5G;addressid/0;aid/5fe6191bf39a42c9;oaid/e3a9473ef6699f75;osVer/29;appBuild/1436;psn/b3rJlGi AwLqa9AqX7Vp0jv4T7XPMa0o|5;psq/4;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/5.4;jdv/;ref/HomeFragment;partner/xiaomi;apprpd/Home_Main;eufv/1;Mozilla/5.0 (Linux; Android 10; Mi9 Pro 5G Build/QKQ1.190825.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045135 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.4;4e6b46913a2e18dd06d6d69843ee4cdd8e033bc1;network/3g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone13,2;addressid/666624049;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/54.11;apprpd/MessageCenter_MessageMerge;ref/MessageCenterController;psq/10;ads/;psn/4e6b46913a2e18dd06d6d69843ee4cdd8e033bc1|101;jdv/0|kong|t_2010804675_|jingfen|810dab1ba2c04b8588c5aa5a0d44c4bd|1614183499;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.2;c71b599e9a0bcbd8d1ad924d85b5715530efad06;network/wifi;ADID/751C6E92-FD10-4323-B37C-187FD0CF0551;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,8;addressid/4053561885;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/263.8;apprpd/;ref/JDLTSubMainPageViewController;psq/2;ads/;psn/c71b599e9a0bcbd8d1ad924d85b5715530efad06|481;jdv/0|kong|t_1001610202_|jingfen|3911bea7ee2f4fcf8d11fdf663192bbe|1614157052210|1614157056;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.2;Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;2d306ee3cacd2c02560627a5113817ebea20a2c9;network/4g;ADID/A346F099-3182-4889-9A62-2B3C28AB861E;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone13,3;hasOCPay/0;appBuild/1017;supportBestPay/0;addressid/;pv/1.35;apprpd/Allowance_Registered;ref/JDLTTaskCenterViewController;psq/0;ads/;psn/2d306ee3cacd2c02560627a5113817ebea20a2c9|2;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;28355aff16cec8bcf3e5728dbbc9725656d8c2c2;network/4g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone10,2;addressid/833058617;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.10;apprpd/;ref/JDLTWebViewController;psq/9;ads/;psn/28355aff16cec8bcf3e5728dbbc9725656d8c2c2|5;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;13.5;24ddac73a3de1b91816b7aedef53e97c4c313733;network/4g;ADID/598C6841-76AC-4512-AA97-CBA940548D70;hasUPPay/0;pushNoticeIsOpen/1;lang/zh_CN;model/iPhone11,6;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/12.6;apprpd/;ref/JDLTSubMainPageViewController;psq/5;ads/;psn/24ddac73a3de1b91816b7aedef53e97c4c313733|23;jdv/0|kong|t_1000170135|tuiguang|notset|1614126110904|1614126110;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 13.5;Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;d7732ba60c8ff73cc3f5ba7290a3aa9551f73a1b;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone12,1;addressid/25239372;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/8.6;apprpd/;ref/JDLTSubMainPageViewController;psq/5;ads/;psn/d7732ba60c8ff73cc3f5ba7290a3aa9551f73a1b|14;jdv/0|kong|t_1001226363_|jingfen|5713234d1e1e4893b92b2de2cb32484d|1614182989528|1614182992;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;ca1a32afca36bc9fb37fd03f18e653bce53eaca5;network/wifi;ADID/3AF380AB-CB74-4FE6-9E7C-967693863CA3;hasUPPay/0;pushNoticeIsOpen/1;lang/zh_CN;model/iPhone8,1;addressid/138323416;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/72.12;apprpd/;ref/JDLTRedPacketViewController;psq/3;ads/;psn/ca1a32afca36bc9fb37fd03f18e653bce53eaca5|109;jdv/0|kong|t_1000536212_|jingfen|c82bfa19e33a4269a5884ffc614790f4|1614141246;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;10;7346933333666353-8333366646039373;network/wifi;model/ONEPLUS A5010;addressid/138117973;aid/7d933f6583cfd097;oaid/;osVer/29;appBuild/1436;psn/T/eqfRSwp8VKEvvXyEunq09Cg2MUkiQ5|17;psq/4;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/11.4;jdv/0|kong|t_1001849073_|jingfen|495a47f6c0b8431c9d460f61ad2304dc|1614084403978|1614084407;ref/HomeFragment;partner/oppo;apprpd/Home_Main;eufv/1;Mozilla/5.0 (Linux; Android 10; ONEPLUS A5010 Build/QKQ1.191014.012; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045140 Mobile Safari/537.36`,
	`jdltapp;android;3.7.0;11;4626269356736353-5353236346334673;network/wifi;model/M2006J10C;addressid/0;aid/dbb9e7655526d3d7;oaid/66a7af49362987b0;osVer/30;appBuild/1436;psn/rQRQgJ 4 S3qkq8YDl28y6jkUHmI/rlX|3;psq/4;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 11;osv/11;pv/3.4;jdv/;ref/HomeFragment;partner/xiaomi;apprpd/Home_Main;eufv/1;Mozilla/5.0 (Linux; Android 11; M2006J10C Build/RP1A.200720.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045513 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.4;78fc1d919de0c8c2de15725eff508d8ab14f9c82;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone13,1;addressid/137829713;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/23.11;apprpd/;ref/JDLTSubMainPageViewController;psq/10;ads/;psn/78fc1d919de0c8c2de15725eff508d8ab14f9c82|34;jdv/0|iosapp|t_335139774|appshare|Wxfriends|1612508702380|1612534293;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;10;0373263343266633-5663030363465326;network/wifi;model/Redmi Note 7;addressid/590846082;aid/07b34bf3e6006d5b;oaid/17975a142e67ec92;osVer/29;appBuild/1436;psn/OHNqtdhQKv1okyh7rB3HxjwI00ixJMNG|4;psq/3;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/2.3;jdv/;ref/activityId=8a8fabf3cccb417f8e691b6774938bc2;partner/xiaomi;apprpd/jsbqd_home;eufv/1;Mozilla/5.0 (Linux; Android 10; Redmi Note 7 Build/QKQ1.190910.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.152 Mobile Safari/537.36`,
	`jdltapp;android;3.7.0;10;3636566623663623-1693635613166646;network/wifi;model/ASUS_I001DA;addressid/1397761133;aid/ccef2fc2a96e1afd;oaid/;osVer/29;appBuild/1436;psn/T8087T0D82PHzJ4VUMGFrfB9dw4gUnKG|76;psq/5;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/73.5;jdv/0|kong|t_1002354188_|jingfen|2335e043b3344107a2750a781fde9a2e#/|1614097081426|1614097087;ref/com.jd.jdlite.lib.personal.view.fragment.JDPersonalFragment;partner/yingyongbao;apprpd/MyJD_Main;eufv/1;Mozilla/5.0 (Linux; Android 10; ASUS_I001DA Build/QKQ1.190825.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045140 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone10,2;addressid/138419019;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/5.7;apprpd/MyJD_Main;ref/MyJdMTAManager;psq/6;ads/;psn/4ee6af0db48fd605adb69b63f00fcbb51c2fc3f0|9;jdv/0|direct|-|none|-|1613705981655|1613823229;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;network/wifi;ADID/F9FD7728-2956-4DD1-8EDD-58B07950864C;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone10,1;addressid/1346909722;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/30.8;apprpd/;ref/JDLTSubMainPageViewController;psq/7;ads/;psn/40d4d4323eb3987226cae367d6b0d8be50f2c7b3|39;jdv/0|kong|t_1000252057_0|tuiguang|eba7648a0f4445aa9cfa6f35c6f36e15|1613995717959|1613995723;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;ADID/5D306F0D-A131-4B26-947E-166CCB9BFFFF;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,6;addressid/138164461;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/7.8;apprpd/;ref/JDLTSubMainPageViewController;psq/7;ads/;psn/d40e5d4a33c100e8527f779557c347569b49c304|7;jdv/0|kong|t_1001226363_|jingfen|3bf5372cb9cd445bbb270b8bc9a34f00|1608439066693|1608439068;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPad;3.7.0;14.5;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPad8,9;hasOCPay/0;appBuild/1017;supportBestPay/0;addressid/;pv/1.20;apprpd/MyJD_Main;ref/MyJdMTAManager;psq/5;ads/;psn/d9f5ddaa0160a20f32fb2c8bfd174fae7993c1b4|3;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.5;Mozilla/5.0 (iPad; CPU OS 14_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;network/wifi;ADID/31548A9C-8A01-469A-B148-E7D841C91FD0;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/10.5;apprpd/;ref/JDLTSubMainPageViewController;psq/4;ads/;psn/a858fb4b40e432ea32f80729916e6c3e910bb922|12;jdv/0|direct|-|none|-|1613898710373|1613898712;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;13.5;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone9,2;addressid/2237496805;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/13.6;apprpd/;ref/JDLTSubMainPageViewController;psq/5;ads/;psn/48e495dcf5dc398b4d46b27e9f15a2b427a154aa|15;jdv/0|direct|-|none|-|1613354874698|1613952828;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 13.5;Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;10;3346332626262353-1666434336539336;network/wifi;model/ONEPLUS A6000;addressid/0;aid/3d3bbb25af44c59c;oaid/;osVer/29;appBuild/1436;psn/ECbc2EqmdSa7mDF1PS1GSrV/Tn7R1LS1|6;psq/8;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/2.67;jdv/0|direct|-|none|-|1613822479379|1613991194;ref/com.jd.jdlite.lib.personal.view.fragment.JDPersonalFragment;partner/oppo;apprpd/MyJD_Main;eufv/1;Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190716.003; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045140 Mobile Safari/537.36`,
	`jdltapp;android;3.7.0;8.1.0;8363834353530333132333132373-43D2930366035323639333662383;network/wifi;model/16th Plus;addressid/0;aid/f909e5f2c464c7c6;oaid/;osVer/27;appBuild/1436;psn/c21YWvVr77Hn6 pOZfxXGY4TZrre1 UOL5hcPbCEDMo=|3;psq/10;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 8.1.0;osv/8.1.0;pv/2.15;jdv/;ref/com.jd.jdlite.lib.personal.view.fragment.JDPersonalFragment;partner/jsxdlyqj09;apprpd/MyJD_Main;eufv/1;Mozilla/5.0 (Linux; Android 8.1.0; 16th Plus Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045514 Mobile Safari/537.36`,
	`jdltapp;android;3.7.0;11;1343467336264693-3343562673463613;network/wifi;model/Mi 10 Pro;addressid/0;aid/14d7cbd934eb7dc1;oaid/335f198546eb3141;osVer/30;appBuild/1436;psn/ZcQh/Wov sNYfZ6JUjTIUBu28 KT0T3u|1;psq/24;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 11;osv/11;pv/1.24;jdv/;ref/com.jd.jdlite.lib.jdlitemessage.view.activity.MessageCenterMainActivity;partner/xiaomi;apprpd/MessageCenter_MessageMerge;eufv/1;Mozilla/5.0 (Linux; Android 11; Mi 10 Pro Build/RKQ1.200826.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.181 Mobile Safari/537.36`,
	`jdltapp;android;3.7.0;10;8353636393732346-6646931673935346;network/wifi;model/MI 8;addressid/1969998059;aid/8566972dfd9a795d;oaid/4a8b773c3e307386;osVer/29;appBuild/1436;psn/PhYbUtCsCJo r 1b8hwxjnY8rEv5S8XC|383;psq/14;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/374.14;jdv/0|iosapp|t_335139774|liteshare|CopyURL|1609306590175|1609306596;ref/com.jd.jdlite.lib.jdlitemessage.view.activity.MessageCenterMainActivity;partner/jsxdlyqj09;apprpd/MessageCenter_MessageMerge;eufv/1;Mozilla/5.0 (Linux; Android 10; MI 8 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045140 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.4;6d343c58764a908d4fa56609da4cb3a5cc1396d3;network/wifi;ADID/4965D884-3E61-4C4E-AEA7-9A8CE3742DA7;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone9,1;addressid/70390480;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.24;apprpd/MyJD_Main;ref/https%3A%2F%2Fjdcs.m.jd.com%2Fafter%2Findex.action%3FcategoryId%3D600%26v%3D6%26entry%3Dm_self_jd;psq/4;ads/;psn/6d343c58764a908d4fa56609da4cb3a5cc1396d3|17;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;13.6.1;4606ddccdfe8f343f8137de7fea7f91fc4aef3a3;network/4g;ADID/C6FB6E20-D334-45FA-818A-7A4C58305202;hasUPPay/0;pushNoticeIsOpen/1;lang/zh_CN;model/iPhone10,1;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/5.9;apprpd/MyJD_Main;ref/MyJdMTAManager;psq/8;ads/;psn/4606ddccdfe8f343f8137de7fea7f91fc4aef3a3|5;jdv/0|iosapp|t_335139774|liteshare|Qqfriends|1614206359106|1614206366;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 13.6.1;Mozilla/5.0 (iPhone; CPU iPhone OS 13_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;3b6e79334551fc6f31952d338b996789d157c4e8;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone10,1;addressid/138051400;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/14.34;apprpd/MyJD_Main;ref/MyJdMTAManager;psq/12;ads/;psn/3b6e79334551fc6f31952d338b996789d157c4e8|46;jdv/0|kong|t_1001707023_|jingfen|e80d7173a4264f4c9a3addcac7da8b5d|1613837384708|1613858760;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;10;1346235693831363-2373837393932673;network/wifi;model/LYA-AL00;addressid/3321567203;aid/1d2e9816278799b7;oaid/00000000-0000-0000-0000-000000000000;osVer/29;appBuild/1436;psn/45VUZFTZJkhP5fAXbeBoQ0   O2GCB I|7;psq/5;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/5.8;jdv/0|iosapp|t_335139774|liteshare|CopyURL|1614066210320|1614066219;ref/com.jd.jdlite.lib.personal.view.fragment.JDPersonalFragment;partner/huawei;apprpd/MyJD_Main;eufv/1;Mozilla/5.0 (Linux; Android 10; LYA-AL00 Build/HUAWEILYA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.106 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.3;c2a8854e622a1b17a6c56c789f832f9d78ef1ba7;network/wifi;hasUPPay/0;pushNoticeIsOpen/1;lang/zh_CN;model/iPhone12,5;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/3.9;apprpd/MyJD_Main;ref/MyJdMTAManager;psq/8;ads/;psn/c2a8854e622a1b17a6c56c789f832f9d78ef1ba7|6;jdv/0|direct|-|none|-|1613541016735|1613823566;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;9;;network/wifi;model/MIX 2S;addressid/;aid/f87efed6d9ed3c65;oaid/94739128ef9dd245;osVer/28;appBuild/1436;psn/R7wD/OWkQjYWxax1pDV6kTIDFPJCUid7C/nl2hHnUuI=|3;psq/13;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 9;osv/9;pv/1.42;jdv/;ref/activityId=8a8fabf3cccb417f8e691b6774938bc2;partner/xiaomi;apprpd/jsbqd_home;eufv/1;Mozilla/5.0 (Linux; Android 9; MIX 2S Build/PKQ1.180729.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.181 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;10;network/wifi;Mozilla/5.0 (Linux; Android 10; Redmi Note 7 Build/QKQ1.190910.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.152 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.4;network/3g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148`,
	`jdltapp;iPad;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/1;lang/zh_CN;model/iPad6,3;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/231.11;pap/JA2020_3112531|3.7.0|IOS 14.4;apprpd/;psn/f5e7b7980fb50efc9c294ac38653c1584846c3db|305;usc/kong;jdv/0|kong|t_1000170135|tuiguang|notset|1613606450668|1613606450;umd/tuiguang;psq/2;ucp/t_1000170135;app_device/IOS;utr/notset;ref/JDLTRedPacketViewController;adk/;ads/;Mozilla/5.0 (iPad; CPU OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;13.5;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone8,1;addressid/669949466;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/9.11;apprpd/;ref/JDLTSubMainPageViewController;psq/10;ads/;psn/500a795cb2abae60b877ee4a1930557a800bef1c|11;jdv/;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 13.5;Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;network/3g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone13,4;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.47;apprpd/;ref/JDLTSubMainPageViewController;psq/8;ads/;psn/21631ed983b3e854a3154b0336413825ad0d6783|9;jdv/0|direct|-|none|-|1614150725100|1614225882;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;network/3g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone13,4;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.47;apprpd/;ref/JDLTSubMainPageViewController;psq/8;ads/;psn/21631ed983b3e854a3154b0336413825ad0d6783|9;jdv/0|direct|-|none|-|1614150725100|1614225882;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone13,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/3.15;apprpd/;ref/https%3A%2F%2Fjdcs.m.jd.com%2Fchat%2Findex.action%3Fentry%3Djd_m_JiSuCommodity%26pid%3D7763388%26lng%3D118.159665%26lat%3D24.504633%26sid%3D31cddc2d58f6e36bf2c31c4e8a79767w%26un_area%3D16_1315_3486_0;psq/12;ads/;psn/c10e0db6f15dec57a94637365f4c3d43e05bbd48|4;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone13,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/3.15;apprpd/;ref/https%3A%2F%2Fjdcs.m.jd.com%2Fchat%2Findex.action%3Fentry%3Djd_m_JiSuCommodity%26pid%3D7763388%26lng%3D118.159665%26lat%3D24.504633%26sid%3D31cddc2d58f6e36bf2c31c4e8a79767w%26un_area%3D16_1315_3486_0;psq/12;ads/;psn/c10e0db6f15dec57a94637365f4c3d43e05bbd48|4;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone13,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/3.15;apprpd/;ref/https%3A%2F%2Fjdcs.m.jd.com%2Fchat%2Findex.action%3Fentry%3Djd_m_JiSuCommodity%26pid%3D7763388%26lng%3D118.159665%26lat%3D24.504633%26sid%3D31cddc2d58f6e36bf2c31c4e8a79767w%26un_area%3D16_1315_3486_0;psq/12;ads/;psn/c10e0db6f15dec57a94637365f4c3d43e05bbd48|4;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,6;hasOCPay/0;appBuild/1017;supportBestPay/0;addressid/2813715704;pv/67.38;apprpd/MyJD_Main;ref/https%3A%2F%2Fh5.m.jd.com%2FbabelDiy%2FZeus%2F2ynE8QDtc2svd36VowmYWBzzDdK6%2Findex.html%3Flng%3D103.957532%26lat%3D30.626962%26sid%3D4fe8ef4283b24723a7bb30ee87c18b2w%26un_area%3D22_1930_49324_52512;psq/4;ads/;psn/5aef178f95931bdbbde849ea9e2fc62b18bc5829|127;jdv/0|direct|-|none|-|1612588090667|1613822580;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;;network/4g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/6.28;apprpd/;ref/JDLTRedPacketViewController;psq/3;ads/;psn/d7beab54ae7758fa896c193b49470204fbb8fce9|8;jdv/0|kong|t_1001707023_|jingfen|79ad0319fa4d47e38521a616d80bc4bd|1613800945610|1613824900;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;network/4g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/6.28;apprpd/;ref/JDLTRedPacketViewController;psq/3;ads/;psn/d7beab54ae7758fa896c193b49470204fbb8fce9|8;jdv/0|kong|t_1001707023_|jingfen|79ad0319fa4d47e38521a616d80bc4bd|1613800945610|1613824900;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;;network/4g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/6.28;apprpd/;ref/JDLTRedPacketViewController;psq/3;ads/;psn/d7beab54ae7758fa896c193b49470204fbb8fce9|8;jdv/0|kong|t_1001707023_|jingfen|79ad0319fa4d47e38521a616d80bc4bd|1613800945610|1613824900;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;network/4g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/6.28;apprpd/;ref/JDLTRedPacketViewController;psq/3;ads/;psn/d7beab54ae7758fa896c193b49470204fbb8fce9|8;jdv/0|kong|t_1001707023_|jingfen|79ad0319fa4d47e38521a616d80bc4bd|1613800945610|1613824900;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;network/4g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,2;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/6.28;apprpd/;ref/JDLTRedPacketViewController;psq/3;ads/;psn/d7beab54ae7758fa896c193b49470204fbb8fce9|8;jdv/0|kong|t_1001707023_|jingfen|79ad0319fa4d47e38521a616d80bc4bd|1613800945610|1613824900;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/4g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone12,1;addressid/3104834020;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.6;apprpd/;ref/JDLTSubMainPageViewController;psq/5;ads/;psn/c633e62b5a4ad0fdd93d9862bdcacfa8f3ecef63|6;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone10,1;addressid/1346909722;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/30.8;apprpd/;ref/JDLTSubMainPageViewController;psq/7;ads/;psn/40d4d4323eb3987226cae367d6b0d8be50f2c7b3|39;jdv/0|kong|t_1000252057_0|tuiguang|eba7648a0f4445aa9cfa6f35c6f36e15|1613995717959|1613995723;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.3;network/wifi;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone10,1;addressid/1346909722;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/30.8;apprpd/;ref/JDLTSubMainPageViewController;psq/7;ads/;psn/40d4d4323eb3987226cae367d6b0d8be50f2c7b3|39;jdv/0|kong|t_1000252057_0|tuiguang|eba7648a0f4445aa9cfa6f35c6f36e15|1613995717959|1613995723;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.3;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,6;addressid/138164461;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/7.8;apprpd/;ref/JDLTSubMainPageViewController;psq/7;ads/;psn/d40e5d4a33c100e8527f779557c347569b49c304|7;jdv/0|kong|t_1001226363_|jingfen|3bf5372cb9cd445bbb270b8bc9a34f00|1608439066693|1608439068;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,6;addressid/138164461;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/7.8;apprpd/;ref/JDLTSubMainPageViewController;psq/7;ads/;psn/d40e5d4a33c100e8527f779557c347569b49c304|7;jdv/0|kong|t_1001226363_|jingfen|3bf5372cb9cd445bbb270b8bc9a34f00|1608439066693|1608439068;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone11,6;addressid/138164461;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/7.8;apprpd/;ref/JDLTSubMainPageViewController;psq/7;ads/;psn/d40e5d4a33c100e8527f779557c347569b49c304|7;jdv/0|kong|t_1001226363_|jingfen|3bf5372cb9cd445bbb270b8bc9a34f00|1608439066693|1608439068;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;13.5;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone9,2;addressid/2237496805;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/13.6;apprpd/;ref/JDLTSubMainPageViewController;psq/5;ads/;psn/48e495dcf5dc398b4d46b27e9f15a2b427a154aa|15;jdv/0|direct|-|none|-|1613354874698|1613952828;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 13.5;Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;android;3.7.0;10;network/wifi;model/ONEPLUS A6000;addressid/0;aid/3d3bbb25af44c59c;oaid/;osVer/29;appBuild/1436;psn/ECbc2EqmdSa7mDF1PS1GSrV/Tn7R1LS1|6;psq/8;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/2.67;jdv/0|direct|-|none|-|1613822479379|1613991194;ref/com.jd.jdlite.lib.personal.view.fragment.JDPersonalFragment;partner/oppo;apprpd/MyJD_Main;eufv/1;Mozilla/5.0 (Linux; Android 10; ONEPLUS A6000 Build/QKQ1.190716.003; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045140 Mobile Safari/537.36`,
	`jdltapp;android;3.7.0;8.1.0;network/wifi;model/16th Plus;addressid/0;aid/f909e5f2c464c7c6;oaid/;osVer/27;appBuild/1436;psn/c21YWvVr77Hn6 pOZfxXGY4TZrre1 UOL5hcPbCEDMo=|3;psq/10;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 8.1.0;osv/8.1.0;pv/2.15;jdv/;ref/com.jd.jdlite.lib.personal.view.fragment.JDPersonalFragment;partner/jsxdlyqj09;apprpd/MyJD_Main;eufv/1;Mozilla/5.0 (Linux; Android 8.1.0; 16th Plus Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045514 Mobile Safari/537.36`,
	`jdltapp;android;3.7.0;11;network/wifi;model/Mi 10 Pro;addressid/0;aid/14d7cbd934eb7dc1;oaid/335f198546eb3141;osVer/30;appBuild/1436;psn/ZcQh/Wov sNYfZ6JUjTIUBu28 KT0T3u|1;psq/24;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 11;osv/11;pv/1.24;jdv/;ref/com.jd.jdlite.lib.jdlitemessage.view.activity.MessageCenterMainActivity;partner/xiaomi;apprpd/MessageCenter_MessageMerge;eufv/1;Mozilla/5.0 (Linux; Android 11; Mi 10 Pro Build/RKQ1.200826.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.181 Mobile Safari/537.36`,
	`jdltapp;android;3.7.0;10;network/wifi;model/MI 8;addressid/1969998059;aid/8566972dfd9a795d;oaid/4a8b773c3e307386;osVer/29;appBuild/1436;psn/PhYbUtCsCJo r 1b8hwxjnY8rEv5S8XC|383;psq/14;adk/;ads/;pap/JA2020_3112531|3.7.0|ANDROID 10;osv/10;pv/374.14;jdv/0|iosapp|t_335139774|liteshare|CopyURL|1609306590175|1609306596;ref/com.jd.jdlite.lib.jdlitemessage.view.activity.MessageCenterMainActivity;partner/jsxdlyqj09;apprpd/MessageCenter_MessageMerge;eufv/1;Mozilla/5.0 (Linux; Android 10; MI 8 Build/QKQ1.190828.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/66.0.3359.126 MQQBrowser/6.2 TBS/045140 Mobile Safari/537.36`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone8,4;addressid/1477231693;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/21.15;apprpd/MyJD_Main;ref/https%3A%2F%2Fgold.jd.com%2F%3Flng%3D0.000000%26lat%3D0.000000%26sid%3D4584eb84dc00141b0d58e000583a338w%26un_area%3D19_1607_3155_62114;psq/0;ads/;psn/2c822e59db319590266cc83b78c4a943783d0077|46;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone9,1;addressid/70390480;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.24;apprpd/MyJD_Main;ref/https%3A%2F%2Fjdcs.m.jd.com%2Fafter%2Findex.action%3FcategoryId%3D600%26v%3D6%26entry%3Dm_self_jd;psq/4;ads/;psn/6d343c58764a908d4fa56609da4cb3a5cc1396d3|17;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone9,1;addressid/70390480;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.24;apprpd/MyJD_Main;ref/https%3A%2F%2Fjdcs.m.jd.com%2Fafter%2Findex.action%3FcategoryId%3D600%26v%3D6%26entry%3Dm_self_jd;psq/4;ads/;psn/6d343c58764a908d4fa56609da4cb3a5cc1396d3|17;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone9,1;addressid/70390480;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.24;apprpd/MyJD_Main;ref/https%3A%2F%2Fjdcs.m.jd.com%2Fafter%2Findex.action%3FcategoryId%3D600%26v%3D6%26entry%3Dm_self_jd;psq/4;ads/;psn/6d343c58764a908d4fa56609da4cb3a5cc1396d3|17;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone9,1;addressid/70390480;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.24;apprpd/MyJD_Main;ref/https%3A%2F%2Fjdcs.m.jd.com%2Fafter%2Findex.action%3FcategoryId%3D600%26v%3D6%26entry%3Dm_self_jd;psq/4;ads/;psn/6d343c58764a908d4fa56609da4cb3a5cc1396d3|17;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPhone;3.7.0;14.4;network/4g;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPhone12,3;hasOCPay/0;appBuild/1017;supportBestPay/0;addressid/;pv/3.49;apprpd/MyJD_Main;ref/MyJdMTAManager;psq/7;ads/;psn/9e0e0ea9c6801dfd53f2e50ffaa7f84c7b40cd15|6;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
	`jdltapp;iPad;3.7.0;14.4;network/wifi;hasUPPay/0;pushNoticeIsOpen/0;lang/zh_CN;model/iPad7,5;addressid/;hasOCPay/0;appBuild/1017;supportBestPay/0;pv/4.14;apprpd/MyJD_Main;ref/MyJdMTAManager;psq/3;ads/;psn/956c074c769cd2eeab2e36fca24ad4c9e469751a|8;jdv/0|;adk/;app_device/IOS;pap/JA2020_3112531|3.7.0|IOS 14.4;Mozilla/5.0 (iPad; CPU OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1`,
}

var ua = func() string {
	return USER_AGENTS[int(time.Now().Unix())%len(USER_AGENTS)]
}

var assets sync.Map
var queryAssetLocker sync.Mutex
var GetAsset = func(ck *JdCookie) string {
	if asset, ok := assets.Load(ck.PtPin); ok {
		return asset.(string)
	}
	queryAssetLocker.Lock()
	defer queryAssetLocker.Unlock()
	var asset = (&JdCookie{
		PtKey: ck.PtKey,
		PtPin: ck.PtPin,
	}).QueryAsset()
	assets.Store(ck.PtPin, asset)
	return asset
}

//检测登录需增加接口
func initAsset() {
	go func() {
		for {
			time.Sleep(time.Minute * 2)
			assets.Range(func(key, _ interface{}) bool {
				assets.Delete(key)
				return true
			})
		}
	}()
	get := func(c chan string, ck JdCookie) {
		c <- GetAsset(&ck)
		return
	}
	//待做：增加惊喜工厂
	core.AddCommand("jd", []core.Function{
		{
			Rules: []string{`asset ?`, `raw ^` + jd_cookie.Get("asset_query_alias", "查询") + ` (\S+)$`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				if s.GetImType() == "wxsv" && !s.IsAdmin() && jd_cookie.GetBool("ban_wxsv") {
					return "不支持此功能。"
				}
				if s.GetImType() == "tg" {
					s.Disappear(time.Second * 40)
				}
				a := s.Get()
				if a == "300" {
					a = "3"
				}
				err, qls := qinglong.QinglongSC(s)
				if err != nil {
					return err
				}
				envs, err := GetEnvs(qls[0], "JD_COOKIE")
				if err != nil {
					return err
				}
				if len(envs) == 0 {
					return "青龙没有京东账号。"
				}
				cks := []JdCookie{}
				for _, env := range envs {
					pt_key := FetchJdCookieValue("pt_key", env.Value)
					pt_pin := FetchJdCookieValue("pt_pin", env.Value)
					if pt_key != "" && pt_pin != "" {
						cks = append(cks, JdCookie{
							PtKey: pt_key,
							PtPin: pt_pin,
							Note:  env.Remarks,
						})
					}
				}
				cks = LimitJdCookie(cks, a)
				if len(cks) == 0 {
					return "没有匹配的京东账号。"
				}
				ke := core.Bucket("wxmp").GetBool("isKe?", false)
				if s.GetImType() == "wxmp" && !ke {
					if len(cks) <= 2 {
						cs := []chan string{}
						for _, ck := range cks {
							c := make(chan string)
							cs = append(cs, c)
							go get(c, ck)
						}
						rt := []string{}
						for _, c := range cs {
							rt = append(rt, <-c)
						}
						s.Reply(strings.Join(rt, "\n\n"))
					} else {
						go func() {
							for _, ck := range cks {
								s.Await(s, func(s core.Sender) interface{} {
									return GetAsset(&ck)
								})
							}
						}()
						return "您有多个账号，输入任意字符将依次为您展示查询结果(公众号查询可能失败，请多试几次)："
					}

				} else {
					for _, ck := range cks {
						tempMsg := GetAsset(&ck)
						translateEmoji(&tempMsg, s.GetImType() == "wx")
						s.Reply(tempMsg)
					}
				}
				return nil
			},
		},
		{
			Rules: []string{`raw ^资产推送$`},
			Cron:  jd_cookie.Get("asset_push"),
			Admin: true,
			Handle: func(_ core.Sender) interface{} {
				qqGroup := jd_cookie.GetInt("qqGroup")
				for _, tp := range []string{
					"qq", "tg", "wx",
				} {
					var fs []func()
					core.Bucket("pin" + strings.ToUpper(tp)).Foreach(func(k, v []byte) error {
						if string(k) != "" {
							jn := &JdNotify{
								ID: string(k),
							}
							jdNotify.First(jn)
							if push, ok := core.Pushs[tp]; ok {
								fs = append(fs, func() {
									push(string(v), GetAsset(&JdCookie{
										PtPin: jn.ID,
										PtKey: jn.PtKey,
									}), qqGroup, "")
								})
							}
						}
						return nil
					})
					if len(fs) != 0 {
						for _, f := range fs {
							f()
						}
					}
					time.Sleep(time.Second)
				}

				return "推送完成"
			},
		},
		{
			Rules: []string{`myCookie`},
			Cron:  jd_cookie.Get("asset_push"),
			Handle: func(s core.Sender) interface{} {
				cookies := []string{}
				tp := s.GetImType()
				uid := s.GetUserID()

				core.Bucket("pin" + strings.ToUpper(tp)).Foreach(func(k, v []byte) error {
					if string(k) != "" && string(v) == uid {
						jn := &JdNotify{
							ID: string(k),
						}
						jdNotify.First(jn)
						cookies = append(cookies, fmt.Sprintf("pt_key=%s;pt_pin=%s;", jn.PtKey, jn.ID))
					}
					return nil
				})

				s.Reply(fmt.Sprintf("已为你找到%d条结果，请在60秒内回复“n”，将依次为你展示。", len(cookies)))
				var ids = []string{}
				for i := range cookies {
					if s.Await(s, func(s core.Sender) interface{} {
						s.RecallMessage(s.GetMessageID())
						return nil
					}, time.Second*60) != "n" {
						return "操作中断。"
					}
					if len(ids) > 0 {
						s.RecallMessage(ids)
					}
					ids, _ = s.Reply(cookies[i])
				}
				if len(ids) > 0 {
					s.RecallMessage(ids)
				}
				return "操作完成。"
			},
		},
		{
			Rules: []string{`^` + jd_cookie.Get("asset_query_alias", "查询") + `$`},
			Handle: func(s core.Sender) interface{} {
				if s.GetImType() == "wxsv" && !s.IsAdmin() && jd_cookie.GetBool("ban_wxsv") {
					return "不支持此功能。"
				}
				if s.GetImType() != "wxmp" {
					go func() {
						l := int64(jd_cookie.GetInt("query_wait_time"))
						if l != 0 {
							deadline := time.Now().Unix() + l
							stop := false
							for {
								if stop {
									break
								}
								s.Await(s, func(_ core.Sender) interface{} {
									left := deadline - time.Now().Unix()
									if left <= 0 {
										stop = true
										left = 1
									}
									return fmt.Sprintf("%d秒后再查询。", left)
								}, "^"+jd_cookie.Get("asset_query_alias", "查询")+"$", time.Second)
							}
						}
					}()
				}
				if groupCode := jd_cookie.Get("groupCode"); !s.IsAdmin() && groupCode != "" && s.GetChatID() != 0 && !strings.Contains(groupCode, fmt.Sprint(s.GetChatID())) {
					return nil
				}
				if query_time := jd_cookie.Get("query_time"); query_time != "" {
					res := regexp.MustCompile(`\d{2}:\d{2}`).FindAllString(query_time, -1)

					if len(res) == 2 {
						n := time.Now().Format("15:04")

						if !(n >= res[0] && n <= res[1]) {
							return query_time
						}
					}
				}
				s.Disappear(time.Second * 40)
				cks := []JdCookie{}
				pin(s.GetImType()).Foreach(func(k, v []byte) error {
					if string(v) == fmt.Sprint(s.GetUserID()) {
						jn := &JdNotify{
							ID: string(k),
						}
						jdNotify.First(jn)
						cks = append(cks, JdCookie{
							PtKey: jn.PtKey,
							PtPin: string(k),
						})
					}
					return nil
				})

				if len(cks) == 0 {
					return "你尚未绑定🐶东账号，请私聊我你的账号信息或者对我说“登录”。"
				}
				if s.GetImType() == "wxmp" {
					cs := []chan string{}
					if len(cks) <= 2 {
						for _, ck := range cks {
							c := make(chan string)
							cs = append(cs, c)
							go get(c, ck)
						}
						rt := []string{}
						for _, c := range cs {
							rt = append(rt, <-c)
						}
						s.Reply(strings.Join(rt, "\n\n"))
					} else {
						go func() {
							for _, ck := range cks {
								s.Await(s, func(s core.Sender) interface{} {
									return GetAsset(&ck)
								})
							}
						}()
						return "您有多个账号，输入任意字符将依次为您展示查询结果："
					}
				} else {
					for _, ck := range cks {
						tempMsg := GetAsset(&ck)
						translateEmoji(&tempMsg, s.GetImType() == "wx")
						s.Reply(tempMsg)
					}
				}
				return nil
			},
		},
		// {
		// 	Rules: []string{`today bean(?)`},
		// 	Admin: true,
		// 	Handle: func(s core.Sender) interface{} {
		// 		a := s.Get()
		// 		envs, err := GetEnvs("JD_COOKIE")
		// 		if err != nil {
		// 			return err
		// 		}
		// 		if len(envs) == 0 {
		// 			return "青龙没有京东账号。"
		// 		}
		// 		cks := []JdCookie{}
		// 		for _, env := range envs {
		// 			pt_key := FetchJdCookieValue("pt_key", env.Value)
		// 			pt_pin := FetchJdCookieValue("pt_pin", env.Value)
		// 			if pt_key != "" && pt_pin != "" {
		// 				cks = append(cks, JdCookie{
		// 					PtKey: pt_key,
		// 					PtPin: pt_pin,
		// 					Note:  env.Remarks,
		// 				})
		// 			}
		// 		}
		// 		cks = LimitJdCookie(cks, a)
		// 		if len(cks) == 0 {
		// 			return "没有匹配的京东账号。"
		// 		}
		// 		var beans []chan int
		// 		for _, ck := range cks {
		// 			var bean = make(chan int)
		// 			go GetTodayBean(&ck, bean)
		// 			beans = append(beans, bean)
		// 		}
		// 		all := 0
		// 		for i := range beans {
		// 			all += <-beans[i]
		// 		}
		// 		return fmt.Sprintf("今日收入%d京豆。", all)
		// 	},
		// },
		// {
		// 	Rules: []string{`yestoday bean(?)`},
		// 	Admin: true,
		// 	Handle: func(s core.Sender) interface{} {
		// 		a := s.Get()
		// 		envs, err := GetEnvs("JD_COOKIE")
		// 		if err != nil {
		// 			return err
		// 		}
		// 		if len(envs) == 0 {
		// 			return "青龙没有京东账号。"
		// 		}
		// 		cks := []JdCookie{}
		// 		for _, env := range envs {
		// 			pt_key := FetchJdCookieValue("pt_key", env.Value)
		// 			pt_pin := FetchJdCookieValue("pt_pin", env.Value)
		// 			if pt_key != "" && pt_pin != "" {
		// 				cks = append(cks, JdCookie{
		// 					PtKey: pt_key,
		// 					PtPin: pt_pin,
		// 					Note:  env.Remarks,
		// 				})
		// 			}
		// 		}
		// 		cks = LimitJdCookie(cks, a)
		// 		if len(cks) == 0 {
		// 			return "没有匹配的京东账号。"
		// 		}
		// 		var beans []chan int
		// 		for _, ck := range cks {
		// 			var bean = make(chan int)
		// 			go GetYestodayBean(&ck, bean)
		// 			beans = append(beans, bean)
		// 		}
		// 		all := 0
		// 		for i := range beans {
		// 			all += <-beans[i]
		// 		}
		// 		return fmt.Sprintf("昨日收入%d京豆。", all)
		// 	},
		// },
		{
			Rules: []string{`imOf ?`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				rt := ""
				pare := s.Get()
				if r := core.FetchCookieValue("pt_pin", pare); r != "" {
					pare = r
				}
				for _, tp := range []string{
					"qq", "tg", "wx", "wxmp",
				} {
					core.Bucket("pin" + strings.ToUpper(tp)).Foreach(func(k, v []byte) error {
						pt_pin := string(k)
						account := string(v)
						if pt_pin == s.Get() && pt_pin != "" {
							rt += fmt.Sprintf("%s - %s\n", tp, account)
						}
						return nil
					})
				}
				if rt == "" {
					return "空"
				}
				return rt
			},
		},
		// {
		// 	Rules: []string{`bean(?)`},
		// 	Admin: true,
		// 	Handle: func(s core.Sender) interface{} {
		// 		a := s.Get()
		// 		envs, err := GetEnvs("JD_COOKIE")
		// 		if err != nil {
		// 			return err
		// 		}
		// 		if len(envs) == 0 {
		// 			return "青龙没有京东账号。"
		// 		}
		// 		cks := []JdCookie{}
		// 		for _, env := range envs {
		// 			pt_key := FetchJdCookieValue("pt_key", env.Value)
		// 			pt_pin := FetchJdCookieValue("pt_pin", env.Value)
		// 			if pt_key != "" && pt_pin != "" {
		// 				cks = append(cks, JdCookie{
		// 					PtKey: pt_key,
		// 					PtPin: pt_pin,
		// 					Note:  env.Remarks,
		// 				})
		// 			}
		// 		}
		// 		cks = LimitJdCookie(cks, a)
		// 		if len(cks) == 0 {
		// 			return "没有匹配的京东账号。"
		// 		}
		// 		all := 0
		// 		for _, ck := range cks {
		// 			ck.Available()
		// 			all += Int(ck.BeanNum)
		// 		}
		// 		return fmt.Sprintf("总计%d京豆。", all)
		// 	},
		// },
	})
}

func LimitJdCookie(cks []JdCookie, a string) []JdCookie {
	ncks := []JdCookie{}
	if s := strings.Split(a, "-"); len(s) == 2 {
		for i := range cks {
			if i+1 >= Int(s[0]) && i+1 <= Int(s[1]) {
				ncks = append(ncks, cks[i])
			}
		}
	} else if x := regexp.MustCompile(`^[\s\d,]+$`).FindString(a); x != "" {
		xx := regexp.MustCompile(`(\d+)`).FindAllStringSubmatch(a, -1)
		for i := range cks {
			for _, x := range xx {
				if i+1 == Int(x[1]) {
					ncks = append(ncks, cks[i])
				}
			}
		}
	}
	if len(ncks) == 0 {
		a = strings.Replace(a, "", "", -1)
		for i := range cks {
			if strings.Contains(cks[i].Note, a) || strings.Contains(cks[i].Nickname, a) || strings.Contains(cks[i].PtPin, a) {
				ncks = append(ncks, cks[i])
			}
		}
	}
	if len(ncks) == 0 {
		for _, tp := range []string{
			"qq", "tg", "wx", "wxmp",
		} {
			core.Bucket("pin" + strings.ToUpper(tp)).Foreach(func(k, v []byte) error {

				pt_pin := string(k)
				account := string(v)
				// fmt.Println(pt_pin, account)
				for _, ck := range cks {
					// fmt.Println(ck.PtPin, pt_pin)
					if ck.PtPin == pt_pin && account == a {
						ncks = append(ncks, ck)
					}
				}
				return nil
			})
		}
	}
	return ncks
}

type Asset struct {
	Nickname string
	Bean     struct {
		Total       int
		TodayIn     int
		TodayOut    int
		YestodayIn  int
		YestodayOut int
		ToExpire    []int
	}
	RedPacket struct {
		Total      float64
		ToExpire   float64
		ToExpireJd float64
		ToExpireJx float64
		ToExpireJs float64
		ToExpireJk float64
		Jd         float64
		Jx         float64
		Js         float64
		Jk         float64
	}
	Other struct {
		JsCoin   float64
		NcStatus float64
		McStatus float64
	}
}

var Int = func(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

var Float64 = func(s string) float64 {
	i, _ := strconv.ParseFloat(s, 64)
	return i
}

type BeanDetail struct {
	Date         string `json:"date"`
	Amount       string `json:"amount"`
	EventMassage string `json:"eventMassage"`
}

func getJingBeanBalanceDetail(page int, cookie string) []BeanDetail {
	type AutoGenerated struct {
		Code       string       `json:"code"`
		DetailList []BeanDetail `json:"detailList"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com/client.action?functionId=getJingBeanBalanceDetail`)
	req.Header("User-Agent", ua())
	req.Header("Host", "api.m.jd.com")
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Cookie", cookie)

	req.Body(fmt.Sprintf(`body={"pageSize": "20", "page": "%d"}&appid=ld`, page))
	data, err := req.Bytes()
	if err != nil {
		return nil
	}
	json.Unmarshal(data, &a)
	return a.DetailList
}

type RedList struct {
	ActivityName string `json:"activityName"`
	Balance      string `json:"balance"`
	BeginTime    int    `json:"beginTime"`
	DelayRemark  string `json:"delayRemark"`
	Discount     string `json:"discount"`
	EndTime      int    `json:"endTime"`
	HbID         string `json:"hbId"`
	HbState      int    `json:"hbState"`
	IsDelay      bool   `json:"isDelay"`
	OrgLimitStr  string `json:"orgLimitStr"`
}

func redPacket(cookie string, rpc chan []RedList) {
	type UseRedInfo struct {
		Count   int       `json:"count"`
		RedList []RedList `json:"redList"`
	}
	type Data struct {
		AvaiCount      int        `json:"avaiCount"`
		Balance        string     `json:"balance"`
		CountdownTime  string     `json:"countdownTime"`
		ExpiredBalance string     `json:"expiredBalance"`
		ServerCurrTime int        `json:"serverCurrTime"`
		UseRedInfo     UseRedInfo `json:"useRedInfo"`
	}
	type AutoGenerated struct {
		Data    Data   `json:"data"`
		Errcode int    `json:"errcode"`
		Msg     string `json:"msg"`
	}
	a := AutoGenerated{}
	req := httplib.Get(`https://m.jingxi.com/user/info/QueryUserRedEnvelopesV2?type=1&orgFlag=JD_PinGou_New&page=1&cashRedType=1&redBalanceFlag=1&channel=3&_=` + fmt.Sprint(time.Now().Unix()) + `&sceneval=2&g_login_type=1&g_ty=ls`)
	req.Header("User-Agent", ua())
	req.Header("authority", "m.jingxi.com")
	req.Header("sec-ch-ua", "\"Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"")
	req.Header("sec-ch-ua-mobile", "?0")
	req.Header("sec-ch-ua-platform", "\"macOS\"")
	// req.Header("accept", "*/*")
	req.Header("sec-fetch-site", "same-site")
	req.Header("sec-fetch-mode", "no-cors")
	req.Header("sec-fetch-dest", "script")
	req.Header("referer", "https://st.jingxi.com/")
	req.Header("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header("Cookie", cookie)
	req.SetTimeout(time.Second, time.Second)
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	rpc <- a.Data.UseRedInfo.RedList
}

func initFarm(cookie string, state chan string) {
	type RightUpResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type TurntableInit struct {
		TimeState int `json:"timeState"`
	}
	type MengchongResouce struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type GUIDPopupTask struct {
		GUIDPopupTask string `json:"guidPopupTask"`
	}
	type IosConfigResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type TodayGotWaterGoalTask struct {
		CanPop bool `json:"canPop"`
	}
	type LeftUpResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type RightDownResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type FarmUserPro struct {
		TotalEnergy     int    `json:"totalEnergy"`
		TreeState       int    `json:"treeState"`
		CreateTime      int64  `json:"createTime"`
		TreeEnergy      int    `json:"treeEnergy"`
		TreeTotalEnergy int    `json:"treeTotalEnergy"`
		ShareCode       string `json:"shareCode"`
		WinTimes        int    `json:"winTimes"`
		NickName        string `json:"nickName"`
		CouponKey       string `json:"couponKey"`
		CouponID        string `json:"couponId"`
		CouponEndTime   int64  `json:"couponEndTime"`
		Type            string `json:"type"`
		SimpleName      string `json:"simpleName"`
		Name            string `json:"name"`
		GoodsImage      string `json:"goodsImage"`
		SkuID           string `json:"skuId"`
		LastLoginDate   int64  `json:"lastLoginDate"`
		NewOldState     int    `json:"newOldState"`
		OldMarkComplete int    `json:"oldMarkComplete"`
		CommonState     int    `json:"commonState"`
		PrizeLevel      int    `json:"prizeLevel"`
	}
	type LeftDownResouces struct {
		AdvertID string `json:"advertId"`
		Name     string `json:"name"`
		AppImage string `json:"appImage"`
		AppLink  string `json:"appLink"`
		CxyImage string `json:"cxyImage"`
		CxyLink  string `json:"cxyLink"`
		Type     string `json:"type"`
		OpenLink bool   `json:"openLink"`
	}
	type LoadFriend struct {
		Code            string      `json:"code"`
		StatisticsTimes interface{} `json:"statisticsTimes"`
		SysTime         int64       `json:"sysTime"`
		Message         interface{} `json:"message"`
		FirstAddUser    bool        `json:"firstAddUser"`
	}
	type AutoGenerated struct {
		Code                  string                `json:"code"`
		RightUpResouces       RightUpResouces       `json:"rightUpResouces"`
		TurntableInit         TurntableInit         `json:"turntableInit"`
		IosShieldConfig       interface{}           `json:"iosShieldConfig"`
		MengchongResouce      MengchongResouce      `json:"mengchongResouce"`
		ClockInGotWater       bool                  `json:"clockInGotWater"`
		GUIDPopupTask         GUIDPopupTask         `json:"guidPopupTask"`
		ToFruitEnergy         int                   `json:"toFruitEnergy"`
		StatisticsTimes       interface{}           `json:"statisticsTimes"`
		SysTime               int64                 `json:"sysTime"`
		CanHongbaoContineUse  bool                  `json:"canHongbaoContineUse"`
		ToFlowTimes           int                   `json:"toFlowTimes"`
		IosConfigResouces     IosConfigResouces     `json:"iosConfigResouces"`
		TodayGotWaterGoalTask TodayGotWaterGoalTask `json:"todayGotWaterGoalTask"`
		LeftUpResouces        LeftUpResouces        `json:"leftUpResouces"`
		MinSupportAPPVersion  string                `json:"minSupportAPPVersion"`
		LowFreqStatus         int                   `json:"lowFreqStatus"`
		FunCollectionHasLimit bool                  `json:"funCollectionHasLimit"`
		Message               interface{}           `json:"message"`
		TreeState             int                   `json:"treeState"`
		RightDownResouces     RightDownResouces     `json:"rightDownResouces"`
		IconFirstPurchaseInit bool                  `json:"iconFirstPurchaseInit"`
		ToFlowEnergy          int                   `json:"toFlowEnergy"`
		FarmUserPro           FarmUserPro           `json:"farmUserPro"`
		RetainPopupLimit      int                   `json:"retainPopupLimit"`
		ToBeginEnergy         int                   `json:"toBeginEnergy"`
		LeftDownResouces      LeftDownResouces      `json:"leftDownResouces"`
		EnableSign            bool                  `json:"enableSign"`
		LoadFriend            LoadFriend            `json:"loadFriend"`
		HadCompleteXgTask     bool                  `json:"hadCompleteXgTask"`
		OldUserIntervalTimes  []int                 `json:"oldUserIntervalTimes"`
		ToFruitTimes          int                   `json:"toFruitTimes"`
		OldUserSendWater      []string              `json:"oldUserSendWater"`
	}
	a := AutoGenerated{}
	req := httplib.Get(`https://api.m.jd.com/client.action?functionId=initForFarm&appid=wh5&body=` + url.QueryEscape(`{"version":4,"channel":1,"babelChannel":"120"}`))
	req.Header("cookie", cookie)
	req.Header("User-Agent", ua())
	req.SetTimeout(time.Second*10, time.Second*10)
	if Transport != nil {
		req.SetTransport(Transport)
	}
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	pt_pin := core.FetchCookieValue("pt_pin", cookie)
	rt := a.FarmUserPro.Name
	not := ""
	if rt == "" {
		rt = "数据异常"
	} else {
		if a.TreeState == 2 || a.TreeState == 3 {
			rt += "已可领取⏰"
			not = rt
		} else if a.TreeState == 1 {
			rt += fmt.Sprintf("种植中，进度%.2f%%🍒", 100*float64(a.FarmUserPro.TreeEnergy)/float64(a.FarmUserPro.TreeTotalEnergy))
		} else if a.TreeState == 0 {
			rt = "您忘了种植新的水果⏰"
			not = rt
		}
	}
	if state != nil {
		state <- rt
	} else if not != "" {
		a叉哦叉哦(pt_pin, "东东农场", not)
	}
}

func initPetTown(cookie string, state chan string) {
	type ResourceList struct {
		AdvertID string `json:"advertId"`
		ImageURL string `json:"imageUrl"`
		Link     string `json:"link"`
		ShopID   string `json:"shopId"`
	}
	type PetPlaceInfoList struct {
		Place  int `json:"place"`
		Energy int `json:"energy"`
	}
	type PetInfo struct {
		AdvertID     string `json:"advertId"`
		NickName     string `json:"nickName"`
		IconURL      string `json:"iconUrl"`
		ClickIconURL string `json:"clickIconUrl"`
		FeedGifURL   string `json:"feedGifUrl"`
		HomePetImage string `json:"homePetImage"`
		CrossBallURL string `json:"crossBallUrl"`
		RunURL       string `json:"runUrl"`
		TickleURL    string `json:"tickleUrl"`
	}
	type GoodsInfo struct {
		GoodsName        string `json:"goodsName"`
		GoodsURL         string `json:"goodsUrl"`
		GoodsID          string `json:"goodsId"`
		ExchangeMedalNum int    `json:"exchangeMedalNum"`
		ActivityID       string `json:"activityId"`
		ActivityIds      string `json:"activityIds"`
	}
	type Result struct {
		ShareCode              string             `json:"shareCode"`
		HisHbFlag              bool               `json:"hisHbFlag"`
		MasterHelpPeoples      []interface{}      `json:"masterHelpPeoples"`
		HelpSwitchOn           bool               `json:"helpSwitchOn"`
		UserStatus             int                `json:"userStatus"`
		TotalEnergy            int                `json:"totalEnergy"`
		MasterInvitePeoples    []interface{}      `json:"masterInvitePeoples"`
		ShareTo                string             `json:"shareTo"`
		PetSportStatus         int                `json:"petSportStatus"`
		UserImage              string             `json:"userImage"`
		MasterHelpReward       int                `json:"masterHelpReward"`
		ShowHongBaoExchangePop bool               `json:"showHongBaoExchangePop"`
		ShowNeedCollectPop     bool               `json:"showNeedCollectPop"`
		PetSportReward         string             `json:"petSportReward"`
		NewhandBubble          bool               `json:"newhandBubble"`
		ResourceList           []ResourceList     `json:"resourceList"`
		ProjectBubble          bool               `json:"projectBubble"`
		MasterInvitePop        bool               `json:"masterInvitePop"`
		MasterInviteReward     int                `json:"masterInviteReward"`
		MedalNum               int                `json:"medalNum"`
		MasterHelpPop          bool               `json:"masterHelpPop"`
		MeetDays               int                `json:"meetDays"`
		PetPlaceInfoList       []PetPlaceInfoList `json:"petPlaceInfoList"`
		MedalPercent           float64            `json:"medalPercent"`
		CharitableSwitchOn     bool               `json:"charitableSwitchOn"`
		PetInfo                PetInfo            `json:"petInfo"`
		NeedCollectEnergy      int                `json:"needCollectEnergy"`
		FoodAmount             int                `json:"foodAmount"`
		InviteCode             string             `json:"inviteCode"`
		RulesURL               string             `json:"rulesUrl"`
		PetStatus              int                `json:"petStatus"`
		GoodsInfo              GoodsInfo          `json:"goodsInfo"`
	}
	type AutoGenerated struct {
		Code       string `json:"code"`
		ResultCode string `json:"resultCode"`
		Message    string `json:"message"`
		Result     Result `json:"result"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com/client.action?functionId=initPetTown`)
	req.Header("Host", "api.m.jd.com")
	req.Header("User-Agent", ua())
	req.Header("cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")

	req.Body(`body={}&appid=wh5&loginWQBiz=pet-town&clientVersion=9.0.4`)
	if Transport != nil {
		req.SetTransport(Transport)
	}
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	rt := ""
	pt_pin := core.FetchCookieValue("pt_pin", cookie)
	not := ""
	if a.Code == "0" && a.ResultCode == "0" && a.Message == "success" {
		if a.Result.UserStatus == 0 {
			rt = "请手动开启活动⏰"
			not = rt

		} else if a.Result.GoodsInfo.GoodsName == "" {
			rt = "你忘了选购新的商品⏰"
			not = rt

		} else if a.Result.PetStatus == 5 {
			rt = a.Result.GoodsInfo.GoodsName + "已可领取⏰"
			not = rt

		} else if a.Result.PetStatus == 6 {
			rt = a.Result.GoodsInfo.GoodsName + "未继续领养新的物品⏰"
			not = rt
		} else {
			rt = a.Result.GoodsInfo.GoodsName + fmt.Sprintf("领养中，进度%.2f%%，勋章%d/%d🐶", a.Result.MedalPercent, a.Result.MedalNum, a.Result.GoodsInfo.ExchangeMedalNum)
		}
	} else {
		rt = "数据异常"
	}
	if state != nil {
		state <- rt
	} else if not != "" {
		a叉哦叉哦(pt_pin, "东东萌宠", not)
	}
}

func jsGold(cookie string, state chan int64) { //

	type BalanceVO struct {
		CashBalance       string `json:"cashBalance"`
		EstimatedAmount   string `json:"estimatedAmount"`
		ExchangeGold      string `json:"exchangeGold"`
		FormatGoldBalance string `json:"formatGoldBalance"`
		GoldBalance       int    `json:"goldBalance"`
	}
	type Gears struct {
		Amount         string `json:"amount"`
		ExchangeAmount string `json:"exchangeAmount"`
		Order          int    `json:"order"`
		Status         int    `json:"status"`
		Type           int    `json:"type"`
	}
	type Data struct {
		Advertise      string    `json:"advertise"`
		BalanceVO      BalanceVO `json:"balanceVO"`
		Gears          []Gears   `json:"gears"`
		IsGetCoupon    bool      `json:"isGetCoupon"`
		IsGetCouponEid bool      `json:"isGetCouponEid"`
		IsLogin        bool      `json:"isLogin"`
		NewPeople      bool      `json:"newPeople"`
	}
	type AutoGenerated struct {
		Code      int    `json:"code"`
		Data      Data   `json:"data"`
		IsSuccess bool   `json:"isSuccess"`
		Message   string `json:"message"`
		RequestID string `json:"requestId"`
	}
	a := AutoGenerated{}
	req := httplib.Post(`https://api.m.jd.com?functionId=MyAssetsService.execute&appid=market-task-h5`)
	req.Header("Accept", "application/json, text/plain, */*")
	// req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Origin", "https://gold.jd.com")
	req.Header("Host", "api.m.jd.com")
	req.Header("Connection", "keep-alive")
	req.Header("User-Agent", ua())
	req.Header("Referer", "https://gold.jd.com/")

	req.Body(`functionId=MyAssetsService.execute&body={"method":"goldShopPage","data":{"channel":1}}&_t=` + fmt.Sprint(time.Now().Unix()) + `&appid=market-task-h5;`)
	data, _ := req.Bytes()
	json.Unmarshal(data, &a)
	if state != nil {
		state <- int64(a.Data.BalanceVO.GoldBalance)
	}
}

func jxncEgg(cookie string, state chan int64) {
	req := httplib.Get("https://m.jingxi.com/jxmc/queryservice/GetHomePageInfo?channel=7&sceneid=1001&activeid=null&activekey=null&isgift=1&isquerypicksite=1&_stk=activeid%2Cactivekey%2Cchannel%2Cisgift%2Cisquerypicksite%2Csceneid&_ste=1&h5st=20210818211830955%3B4408816258824161%3B10028%3Btk01w8db21b2130ny2eg0siAPpNQgBqjGzYfuG6IP7Z%2BAOB40BiqLQ%2Blglfi540AB%2FaQrTduHbnk61ngEeKn813gFeRD%3Bd9a0b833bf99a29ed726cbffa07ba955cc27d1ff7d2d55552878fc18fc667929&_=1629292710957&sceneval=2&g_login_type=1&g_ty=ls")
	req.Header("User-Agent", ua())
	req.Header("Host", "m.jingxi.com")
	// req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	// req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://st.jingxi.com/pingou/jxmc/index.html?nativeConfig=%7B%22immersion%22%3A1%2C%22toColor%22%3A%22%23e62e0f%22%7D&;__mcwvt=sjcp&ptag=7155.9.95")

	req.Header("Cookie", cookie)
	data, _ := req.Bytes()

	egg, _ := jsonparser.GetInt(data, "data", "eggcnt")
	state <- egg
}

func tytCoupon(cookie string, state chan string) {

	type DiscountInfo struct {
		High string        `json:"high"`
		Info []interface{} `json:"info"`
	}
	type ExtInfo struct {
		Num5              string `json:"5"`
		Num12             string `json:"12"`
		Num16             string `json:"16"`
		Num21             string `json:"21"`
		Num52             string `json:"52"`
		Num54             string `json:"54"`
		Num74             string `json:"74"`
		BusinessLabel     string `json:"business_label"`
		LimitOrganization string `json:"limit_organization"`
		UserLabel         string `json:"user_label"`
	}
	type Useable struct {
		AreaDesc         string        `json:"areaDesc"`
		AreaType         int           `json:"areaType"`
		Batchid          string        `json:"batchid"`
		BeanNumForPerson int           `json:"beanNumForPerson"`
		BeanNumForPlat   int           `json:"beanNumForPlat"`
		BeginTime        string        `json:"beginTime"`
		CanBeSell        bool          `json:"canBeSell"`
		CanBeShare       bool          `json:"canBeShare"`
		CompleteTime     string        `json:"completeTime"`
		CouponKind       int           `json:"couponKind"`
		CouponStyle      int           `json:"couponStyle"`
		CouponTitle      string        `json:"couponTitle"`
		Couponid         string        `json:"couponid"`
		Coupontype       int           `json:"coupontype"`
		CreateTime       string        `json:"createTime"`
		Discount         string        `json:"discount"`
		DiscountInfo     DiscountInfo  `json:"discountInfo"`
		EndTime          string        `json:"endTime"`
		ExpireType       int           `json:"expireType"`
		ExtInfo          ExtInfo       `json:"extInfo"`
		HourCoupon       int           `json:"hourCoupon"`
		IsOverlay        int           `json:"isOverlay"`
		LimitStr         string        `json:"limitStr"`
		LinkStr          string        `json:"linkStr"`
		OperateTime      string        `json:"operateTime"`
		OrderID          string        `json:"orderId"`
		OverlayDesc      string        `json:"overlayDesc"`
		PassKey          string        `json:"passKey"`
		Pin              string        `json:"pin"`
		PlatFormInfo     string        `json:"platFormInfo"`
		Platform         int           `json:"platform"`
		PlatformDetails  []interface{} `json:"platformDetails"`
		PwdKey           string        `json:"pwdKey"`
		Quota            string        `json:"quota"`
		SellID           string        `json:"sellId"`
		ShareID          string        `json:"shareId"`
		ShopID           string        `json:"shopId"`
		ShopName         string        `json:"shopName"`
		State            int           `json:"state"`
		UseTime          string        `json:"useTime"`
		VenderID         string        `json:"venderId"`
	}
	type Coupon struct {
		Curtimestamp           int       `json:"curtimestamp"`
		ExpiredCount           int       `json:"expired_count"`
		IsHideBaiTiaoInJxWxapp int       `json:"isHideBaiTiaoInJxWxapp"`
		IsHideMailInWxapp      int       `json:"isHideMailInWxapp"`
		Useable                []Useable `json:"useable"`
		UseableCount           int       `json:"useable_count"`
		UsedCount              int       `json:"used_count"`
	}
	type AutoGenerated struct {
		Coupon    Coupon `json:"coupon"`
		ErrMsg    string `json:"errMsg"`
		ErrorCode int    `json:"errorCode"`
		HasNext   int    `json:"hasNext"`
		Jdpin     string `json:"jdpin"`
		State     int    `json:"state"`
		Uin       string `json:"uin"`
	}
	a := AutoGenerated{}
	req := httplib.Get(`https://m.jingxi.com/activeapi/queryjdcouponlistwithfinance?state=1&wxadd=1&filterswitch=1&_=1629296270692&sceneval=2&g_login_type=1&callback=jsonpCBKB&g_ty=ls`)
	// req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	// req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Host", "m.jingxi.com")

	req.Header("User-Agent", ua())
	req.Header("Referer", "https://st.jingxi.com/my/coupon/jx.shtml?sceneval=2&ptag=7155.1.18")

	data, _ := req.Bytes()
	res := regexp.MustCompile(`jsonpCBKB[(](.*)\s+[)];}catch`).FindSubmatch(data)
	rt := ""
	if len(res) > 0 {
		json.Unmarshal(res[1], &a)
		num := 0
		toexp := 0
		tm := int(time.Now().Unix() * 1000)
		for _, cp := range a.Coupon.Useable {
			if strings.Contains(cp.CouponTitle, "推推5.01") {
				num++
				if Int(cp.EndTime) < tm {
					toexp++
				}
			}
		}
		if num == 0 {
			rt = "无"
		} else {
			rt = fmt.Sprintf("%d张5元优惠券(今天过期)", num)
		}
	}
	state <- rt
}

func mmCoin(cookie string, state chan int64) {
	req := httplib.Post(`https://api.m.jd.com/client.action`)
	req.Header("Host", "api.m.jd.com")
	// req.Header("Accept", "application/json, text/plain, */*")
	req.Header("Origin", "https://h5.m.jd.com")

	req.Header("User-Agent", ua())
	req.Header("cookie", cookie)
	req.Header("Content-Type", "application/x-www-form-urlencoded")

	req.Body(`uuid=3245ad3d16ab2153c69f9ca91cd2e931b06a3bb8&clientVersion=10.1.0&client=wh5&osVersion=&area=&networkType=wifi&functionId=homePageV2&body=%7B%7D&appid=SecKill2020`)
	data, _ := req.Bytes()
	mmc, _ := jsonparser.GetInt(data, "result", "assignment", "assignmentPoints")
	state <- mmc
}

func jdzz(cookie string, state chan int64) { //
	req := httplib.Post(`https://api.m.jd.com/client.action?functionId=interactTaskIndex&body={}&client=wh5&clientVersion=9.1.0`)
	req.Header("Host", "api.m.jd.com")
	req.Header("Accept-Language", "zh-cn")
	// req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "http://wq.jd.com/wxapp/pages/hd-interaction/index/index")
	req.Header("User-Agent", ua())
	req.Header("cookie", cookie)
	req.Header("Content-Type", "application/json")

	data, _ := req.Bytes()
	mmc, _ := jsonparser.GetString(data, "data", "totalNum")
	state <- int64(Int(mmc))
}

func (ck *JdCookie) Available() bool {
	if ck.PtKey == "" {
		return false
	}
	cookie := "pt_key=" + ck.PtKey + ";pt_pin=" + ck.PtPin + ";"
	if ck == nil {
		return true
	}
	req := httplib.Get("https://me-api.jd.com/user_new/info/GetJDUserInfoUnion?orgFlag=JD_PinGou_New&callSource=mainorder&channel=4&isHomewhite=0&sceneval=2&_=" + fmt.Sprint(time.Now().Unix()) + "&sceneval=2&g_login_type=1&g_ty=ls")
	req.Header("Cookie", cookie)
	req.Header("authority", "me-api.jd.com")
	// req.Header("accept", "*/*")
	req.Header("sec-fetch-site", "same-site")
	req.Header("sec-fetch-mode", "no-cors")
	req.Header("sec-fetch-dest", "script")
	req.Header("referer", "https://home.m.jd.com/")
	req.Header("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header("User-Agent", ua())

	data, err := req.Bytes()
	if err != nil {
		return av12(ck)
	}
	ui := &UserInfoResult{}
	if nil != json.Unmarshal(data, ui) {
		return av12(ck)
	}
	switch ui.Retcode {
	// case "1001": //ck.BeanNum
	// 	if ui.Msg == "not login"{
	// 		return false
	// 	}
	case "0":
		realPin := url.QueryEscape(ui.Data.UserInfo.BaseInfo.CurPin)
		if realPin != ck.PtPin {
			if realPin == "" {
				return av12(ck)
			} else {
				ck.PtPin = realPin
			}
		}
		if ui.Data.UserInfo.BaseInfo.Nickname != ck.Nickname || ui.Data.AssetInfo.BeanNum != ck.BeanNum || ui.Data.UserInfo.BaseInfo.UserLevel != ck.UserLevel || ui.Data.UserInfo.BaseInfo.LevelName != ck.LevelName {
			ck.UserLevel = ui.Data.UserInfo.BaseInfo.UserLevel
			ck.LevelName = ui.Data.UserInfo.BaseInfo.LevelName
			ck.Nickname = ui.Data.UserInfo.BaseInfo.Nickname
			ck.BeanNum = ui.Data.AssetInfo.BeanNum
		}
		return true
	}
	return av12(ck)
}

func av12(ck *JdCookie) bool {
	type AutoGenerated struct {
		Sid           string `json:"sid"`
		DesPin        string `json:"desPin"`
		UserFlagCheck string `json:"userFlagCheck"`
		User          struct {
			Balance               interface{} `json:"balance"`
			Birthday              string      `json:"birthday"`
			City                  int         `json:"city"`
			Code                  interface{} `json:"code"`
			Companys              string      `json:"companys"`
			County                int         `json:"county"`
			Coupon                interface{} `json:"coupon"`
			Email                 string      `json:"email"`
			EmployeeInfo          interface{} `json:"employeeInfo"`
			GiftCard              interface{} `json:"giftCard"`
			GiftCardPlusgiftECard string      `json:"giftCardPlusgiftECard"`
			GiftECard             interface{} `json:"giftECard"`
			HomePage              string      `json:"homePage"`
			ImgFlag               int         `json:"imgFlag"`
			ImgURL                string      `json:"imgUrl"`
			IPAddress             string      `json:"ipAddress"`
			JingBean              string      `json:"jingBean"`
			Labels                interface{} `json:"labels"`
			LastTime              string      `json:"lastTime"`
			MiddleSchool          string      `json:"middleSchool"`
			Msn                   string      `json:"msn"`
			MyJdNavigation        interface{} `json:"myJdNavigation"`
			PetName               string      `json:"petName"`
			PlusHeadBkImg         interface{} `json:"plusHeadBkImg"`
			PlusText              interface{} `json:"plusText"`
			Province              int         `json:"province"`
			QianbaoDegradeStatus  bool        `json:"qianbaoDegradeStatus"`
			Qq                    string      `json:"qq"`
			RegIP                 string      `json:"regIp"`
			RegTime               string      `json:"regTime"`
			Remark                string      `json:"remark"`
			SchoolID              int         `json:"schoolId"`
			SchoolYn              int         `json:"schoolYn"`
			Score                 int         `json:"score"`
			SecoSchool            string      `json:"secoSchool"`
			Sex                   int         `json:"sex"`
			ShopIntegral          interface{} `json:"shopIntegral"`
			Uclass                string      `json:"uclass"`
			UnColleger            string      `json:"unColleger"`
			UnickName             string      `json:"unickName"`
			UserID                int         `json:"userId"`
			UserName              string      `json:"userName"`
			UserPlusStatus        bool        `json:"userPlusStatus"`
		} `json:"user"`
	}
	req := httplib.Get(`https://wxapp.m.jd.com/kwxhome/myJd/home.json?&useGuideModule=0&bizId=&brandId=&fromType=wxapp&timestamp=` + fmt.Sprint(time.Now().Unix()))
	req.Header("User-Agent", ua())
	req.Header("Cookie", "pt_key="+ck.PtKey+";pt_pin="+ck.PtPin+";")
	data, err := req.Bytes()
	if err != nil {
		return av2(ck)
	}
	code, _ := jsonparser.GetString(data, "code")
	if code == "999" {
		return false
	} else {
		a := AutoGenerated{}
		if nil != json.Unmarshal(data, &a) {
			return av3(ck)
		}
		if a.User.UnickName == "" {
			return av3(ck)
		}
		ck.Nickname = a.User.UnickName
		ck.BeanNum = a.User.JingBean
		return true
	}

}

func av2(ck *JdCookie) bool {
	req := httplib.Get(`https://m.jingxi.com/user/info/GetJDUserBaseInfo?_=1629334995401&sceneval=2&g_login_type=1&g_ty=ls`)
	req.Header("User-Agent", ua())
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Referer", "https://st.jingxi.com/my/userinfo.html?&ptag=7205.12.4")
	req.Header("Cookie", "pt_key="+ck.PtKey+";pt_pin="+ck.PtPin+";")

	data, err := req.Bytes()
	if err != nil {
		return true
	}
	ck.Nickname, _ = jsonparser.GetString(data, "nickname")
	return !strings.Contains(string(data), "login")
}

func av3(ck *JdCookie) bool {
	req := httplib.Get(`https://wq.jd.com/user_new/info/GetJDUserInfoUnion?sceneval=2`)
	req.Header("User-Agent", ua())
	req.Header("Host", "wq.jd.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Referer", "https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&")

	req.Header("Cookie", "pt_key="+ck.PtKey+";pt_pin="+ck.PtPin+";")
	data, err := req.Bytes()
	if err != nil {
		return av2(ck)
	}
	ck.Nickname, _ = jsonparser.GetString(data, "data", "userInfo", "baseInfo", "nickname")
	ck.BeanNum, _ = jsonparser.GetString(data, "data", "assetInfo", "beanNum")
	if ck.Nickname != "" {
		return true
	} else {
		return av2(ck)
	}
}

type UserInfoResult struct {
	Data struct {
		JdVvipCocoonInfo struct {
			JdVvipCocoon struct {
				DisplayType   int    `json:"displayType"`
				HitTypeList   []int  `json:"hitTypeList"`
				Link          string `json:"link"`
				Price         string `json:"price"`
				Qualification int    `json:"qualification"`
				SellingPoints string `json:"sellingPoints"`
			} `json:"JdVvipCocoon"`
			JdVvipCocoonStatus string `json:"JdVvipCocoonStatus"`
		} `json:"JdVvipCocoonInfo"`
		JdVvipInfo struct {
			JdVvipStatus string `json:"jdVvipStatus"`
		} `json:"JdVvipInfo"`
		AssetInfo struct {
			AccountBalance string `json:"accountBalance"`
			BaitiaoInfo    struct {
				AvailableLimit     string `json:"availableLimit"`
				BaiTiaoStatus      string `json:"baiTiaoStatus"`
				Bill               string `json:"bill"`
				BillOverStatus     string `json:"billOverStatus"`
				Outstanding7Amount string `json:"outstanding7Amount"`
				OverDueAmount      string `json:"overDueAmount"`
				OverDueCount       string `json:"overDueCount"`
				UnpaidForAll       string `json:"unpaidForAll"`
				UnpaidForMonth     string `json:"unpaidForMonth"`
			} `json:"baitiaoInfo"`
			BeanNum    string `json:"beanNum"`
			CouponNum  string `json:"couponNum"`
			CouponRed  string `json:"couponRed"`
			RedBalance string `json:"redBalance"`
		} `json:"assetInfo"`
		FavInfo struct {
			FavDpNum    string `json:"favDpNum"`
			FavGoodsNum string `json:"favGoodsNum"`
			FavShopNum  string `json:"favShopNum"`
			FootNum     string `json:"footNum"`
			IsGoodsRed  string `json:"isGoodsRed"`
			IsShopRed   string `json:"isShopRed"`
		} `json:"favInfo"`
		GrowHelperCoupon struct {
			AddDays     int     `json:"addDays"`
			BatchID     int     `json:"batchId"`
			CouponKind  int     `json:"couponKind"`
			CouponModel int     `json:"couponModel"`
			CouponStyle int     `json:"couponStyle"`
			CouponType  int     `json:"couponType"`
			Discount    float64 `json:"discount"`
			LimitType   int     `json:"limitType"`
			MsgType     int     `json:"msgType"`
			Quota       float64 `json:"quota"`
			RoleID      int     `json:"roleId"`
			State       int     `json:"state"`
			Status      int     `json:"status"`
		} `json:"growHelperCoupon"`
		KplInfo struct {
			KplInfoStatus string `json:"kplInfoStatus"`
			Mopenbp17     string `json:"mopenbp17"`
			Mopenbp22     string `json:"mopenbp22"`
		} `json:"kplInfo"`
		OrderInfo struct {
			CommentCount     string        `json:"commentCount"`
			Logistics        []interface{} `json:"logistics"`
			OrderCountStatus string        `json:"orderCountStatus"`
			ReceiveCount     string        `json:"receiveCount"`
			WaitPayCount     string        `json:"waitPayCount"`
		} `json:"orderInfo"`
		PlusPromotion struct {
			Status int `json:"status"`
		} `json:"plusPromotion"`
		UserInfo struct {
			BaseInfo struct {
				AccountType    string `json:"accountType"`
				BaseInfoStatus string `json:"baseInfoStatus"`
				CurPin         string `json:"curPin"`
				DefinePin      string `json:"definePin"`
				HeadImageURL   string `json:"headImageUrl"`
				LevelName      string `json:"levelName"`
				Nickname       string `json:"nickname"`
				Pinlist        string `json:"pinlist"`
				UserLevel      string `json:"userLevel"`
			} `json:"baseInfo"`
			IsHideNavi     string `json:"isHideNavi"`
			IsHomeWhite    string `json:"isHomeWhite"`
			IsJTH          string `json:"isJTH"`
			IsKaiPu        string `json:"isKaiPu"`
			IsPlusVip      string `json:"isPlusVip"`
			IsQQFans       string `json:"isQQFans"`
			IsRealNameAuth string `json:"isRealNameAuth"`
			IsWxFans       string `json:"isWxFans"`
			Jvalue         string `json:"jvalue"`
			OrderFlag      string `json:"orderFlag"`
			PlusInfo       struct {
			} `json:"plusInfo"`
			XbScore string `json:"xbScore"`
		} `json:"userInfo"`
		UserLifeCycle struct {
			IdentityID      string `json:"identityId"`
			LifeCycleStatus string `json:"lifeCycleStatus"`
			TrackID         string `json:"trackId"`
		} `json:"userLifeCycle"`
	} `json:"data"`
	Msg       string `json:"msg"`
	Retcode   string `json:"retcode"`
	Timestamp int64  `json:"timestamp"`
}

func FetchJdCookieValue(ps ...string) string {
	var key, cookies string
	if len(ps) == 2 {
		if len(ps[0]) > len(ps[1]) {
			key, cookies = ps[1], ps[0]
		} else {
			key, cookies = ps[0], ps[1]
		}
	}
	match := regexp.MustCompile(key + `=([^;]*);{0,1}`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}

func GetTodayBean(ck *JdCookie, state chan int) {
	cookie := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
	today := time.Now().Local().Format("2006-01-02")
	page := 1
	end := false
	in := 0
	defer func() {
		state <- in
	}()
	for {
		if end {
			return
		}
		bds := getJingBeanBalanceDetail(page, cookie)
		if bds == nil {
			break
		}
		for _, bd := range bds {
			amount := Int(bd.Amount)
			if strings.Contains(bd.Date, today) {
				if amount > 0 {
					in += amount
				} else {

				}
			} else {
				end = true
				break
			}
		}
		page++
	}
	return
}

func GetYestodayBean(ck *JdCookie, state chan int) {
	cookie := fmt.Sprintf("pt_key=%s;pt_pin=%s;", ck.PtKey, ck.PtPin)
	today := time.Now().Local().Format("2006-01-02")
	yestoday := time.Now().Local().Add(-time.Hour * 24).Format("2006-01-02")
	page := 1
	end := false
	in := 0
	defer func() {
		state <- in
	}()
	for {
		if end {
			return
		}
		bds := getJingBeanBalanceDetail(page, cookie)
		if bds == nil {
			break
		}
		for _, bd := range bds {
			amount := Int(bd.Amount)
			if strings.Contains(bd.Date, yestoday) {
				if amount > 0 {
					in += amount
				} else {

				}
			} else if strings.Contains(bd.Date, today) {

			} else {
				end = true
				break
			}
		}
		page++
	}
	return
}

type XBeanDetail struct {
	Amount      int    `json:"amount"`
	Createdate  string `json:"createdate"`
	Visibleinfo string `json:"visibleinfo"`
}

//欢迎叼毛来抄代码
func queryuserjingdoudetail(cookie string, e下水道 chan []int) {
	type AutoGenerated struct {
		Detail []XBeanDetail `json:"detail"`
		Ret    int           `json:"ret"`
		Retmsg string        `json:"retmsg"`
	}
	a := AutoGenerated{}
	req := httplib.Get(`https://m.jingxi.com/activeapi/queryuserjingdoudetail?pagesize=10&type=16&_=` + fmt.Sprint(time.Now().Unix()) + `&sceneval=2&g_login_type=1&g_ty=ls`)
	req.Header("User-Agent", ua())
	req.Header("Cookie", cookie)
	req.Header("authority", "m.jingxi.com")
	req.Header("accept", "*/*")
	req.Header("sec-fetch-site", "same-site")
	req.Header("sec-fetch-mode", "no-cors")
	req.Header("sec-fetch-dest", "script")
	req.Header("referer", "https://st.jingxi.com/")
	req.Header("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")

	data, _ := req.Bytes()

	json.Unmarshal(data, &a)
	e叼毛 := []int{0, 0, 0, 0, 0}
	today := time.Now().Local().Format("2006/01/02")
	yestoday := time.Now().Local().Add(-time.Hour * 24).Format("2006/01/02")
	for _, v := range a.Detail {
		e叼毛[0] += v.Amount
		if strings.Contains(v.Createdate, today) {
			if v.Amount > 0 {
				e叼毛[1] += v.Amount
			} else {
				e叼毛[2] += -v.Amount
			}
		} else if strings.Contains(v.Createdate, yestoday) {
			if v.Amount > 0 {
				e叼毛[3] += v.Amount
			} else {
				e叼毛[4] += -v.Amount
			}
		}
	}
	e下水道 <- e叼毛 //叼毛去下水道
}

//欢迎叼毛来抄代码
func jingtie(cookie string, e下水道 chan string) {
	req := httplib.Post(`https://ms.jr.jd.com/gw/generic/jrm/h5/m/channelUserSubsidyInfo`)
	req.Header("User-Agent", ua())
	req.Header("Cookie", cookie)
	req.Body(`"reqData=%7B%22source%22%3A%22H5%22%2C%22channel%22%3A%22default%22%2C%22channelLv%22%3A%22%22%2C%22apiVersion%22%3A%224.0.0%22%2C%22riskDeviceParam%22%3A%22%7B%5C%22macAddress%5C%22%3A%5C%22%5C%22%2C%5C%22imei%5C%22%3A%5C%22%5C%22%2C%5C%22eid%5C%22%3A%5C%22%5C%22%2C%5C%22openUUID%5C%22%3A%5C%22%5C%22%2C%5C%22uuid%5C%22%3A%5C%22%5C%22%2C%5C%22traceIp%5C%22%3A%5C%22%5C%22%2C%5C%22os%5C%22%3A%5C%22%5C%22%2C%5C%22osVersion%5C%22%3A%5C%22%5C%22%2C%5C%22appId%5C%22%3A%5C%22%5C%22%2C%5C%22clientVersion%5C%22%3A%5C%22%5C%22%2C%5C%22resolution%5C%22%3A%5C%22%5C%22%2C%5C%22channelInfo%5C%22%3A%5C%22%5C%22%2C%5C%22networkType%5C%22%3A%5C%22%5C%22%2C%5C%22startNo%5C%22%3A42%2C%5C%22openid%5C%22%3A%5C%22%5C%22%2C%5C%22token%5C%22%3A%5C%22%5C%22%2C%5C%22sid%5C%22%3A%5C%22%5C%22%2C%5C%22terminalType%5C%22%3A%5C%22%5C%22%2C%5C%22longtitude%5C%22%3A%5C%22%5C%22%2C%5C%22latitude%5C%22%3A%5C%22%5C%22%2C%5C%22securityData%5C%22%3A%5C%22%5C%22%2C%5C%22jscContent%5C%22%3A%5C%22%5C%22%2C%5C%22fnHttpHead%5C%22%3A%5C%22%5C%22%2C%5C%22receiveRequestTime%5C%22%3A%5C%22%5C%22%2C%5C%22port%5C%22%3A80%2C%5C%22appType%5C%22%3A%5C%22%5C%22%2C%5C%22deviceType%5C%22%3A%5C%22%5C%22%2C%5C%22fp%5C%22%3A%5C%22%5C%22%2C%5C%22ip%5C%22%3A%5C%22%5C%22%2C%5C%22idfa%5C%22%3A%5C%22%5C%22%2C%5C%22sdkToken%5C%22%3A%5C%22%5C%22%7D%22%2C%22others%22%3A%7B%22shareId%22%3A%22%22%7D%7D"`)
	data, _ := req.String()
	e叼毛 := ""
	res := regexp.MustCompile(`"availableAmount":([^,]+),`).FindStringSubmatch(data)
	if len(res) > 0 {
		e叼毛 = res[1]
	}
	e下水道 <- e叼毛 //叼毛去下水道
}

//欢迎叼毛来抄代码
func jingxiangzhi(cookie string, e下水道 chan string) {
	req := httplib.Get(`https://wxapp.m.jd.com/kwxhome/myJd/home.json?&useGuideModule=0&bizId=&brandId=&fromType=wxapp&timestamp=` + fmt.Sprint(time.Now().Unix()))
	req.Header("User-Agent", ua())
	req.Header("Cookie", cookie)
	data, _ := req.Bytes()
	e叼毛, _ := jsonparser.GetString(data, "user", "uclass")
	e叼毛 = strings.Replace(e叼毛, "京享值", "", -1)
	e下水道 <- e叼毛 //叼毛去下水道
}

//欢迎叼毛来抄代码
func dream(cookie string, state chan string) {
	type AssistCondition struct {
		AssistConditionMsg    string `json:"assistConditionMsg"`
		AssistNumCurrent      int    `json:"assistNumCurrent"`
		AssistNumLimit        int    `json:"assistNumLimit"`
		AssistNumMax          int    `json:"assistNumMax"`
		AssistRemindKey       string `json:"assistRemindKey"`
		AssistRemindUser      string `json:"assistRemindUser"`
		CommodityAppLimitFlag int    `json:"commodityAppLimitFlag"`
		FactoryStatus         int    `json:"factoryStatus"`
		HireNumLimit          int    `json:"hireNumLimit"`
		ReAssistFlag          int    `json:"reAssistFlag"`
		SharePin              string `json:"sharePin"`
		SharePinHeadImage     string `json:"sharePinHeadImage"`
	}
	type AssistMaterialTuanCondition struct {
		AssistAppFlag    int           `json:"assistAppFlag"`
		AssistSelfFlag   int           `json:"assistSelfFlag"`
		CommodityList    []interface{} `json:"commodityList"`
		LimitTime        int           `json:"limitTime"`
		MaterialName     string        `json:"materialName"`
		MaterialPicture  string        `json:"materialPicture"`
		MaterialStatus   int           `json:"materialStatus"`
		OutOfStockFlag   int           `json:"outOfStockFlag"`
		RemindMsg        string        `json:"remindMsg"`
		SharePin         string        `json:"sharePin"`
		SharePinNickname string        `json:"sharePinNickname"`
		StartTime        int           `json:"startTime"`
		TuanID           string        `json:"tuanId"`
	}
	type DeviceList struct {
		CreateTime  int `json:"createTime"`
		DeviceDimID int `json:"deviceDimId"`
		DeviceID    int `json:"deviceId"`
		FactoryID   int `json:"factoryId"`
		UpdateTime  int `json:"updateTime"`
	}
	type FactoryList struct {
		CreateTime int    `json:"createTime"`
		FactoryID  int    `json:"factoryId"`
		Name       string `json:"name"`
		UpdateTime int    `json:"updateTime"`
	}
	type NewFactoryFlower struct {
		FactoryFlowerSendFlag int `json:"factoryFlowerSendFlag"`
		SendElectric          int `json:"sendElectric"`
	}
	type PickSiteInfo struct {
		Address         string `json:"address"`
		CityID          int    `json:"cityId"`
		CityName        string `json:"cityName"`
		CountryID       int    `json:"countryId"`
		CountryName     string `json:"countryName"`
		DcID            int    `json:"dcId"`
		ProvinceID      int    `json:"provinceId"`
		ProvinceName    string `json:"provinceName"`
		Sid             int    `json:"sid"`
		SiteID          string `json:"siteId"`
		SiteName        string `json:"siteName"`
		SiteURL         string `json:"siteUrl"`
		ToastChangeSite bool   `json:"toastChangeSite"`
		TownID          int    `json:"townId"`
		TownName        string `json:"townName"`
		Weight          int    `json:"weight"`
	}
	type ProductionList struct {
		BeginTime        int   `json:"beginTime"`
		CommodityDimID   int   `json:"commodityDimId"`
		CreateTime       int   `json:"createTime"`
		DataMark         int   `json:"dataMark"`
		DeviceID         int   `json:"deviceId"`
		EndTime          int   `json:"endTime"`
		ExchangeStatus   int   `json:"exchangeStatus"`
		FactoryID        int   `json:"factoryId"`
		InvestedElectric int   `json:"investedElectric"`
		NeedElectric     int   `json:"needElectric"`
		ProductionID     int64 `json:"productionId"`
		Status           int   `json:"status"`
		UpdateTime       int   `json:"updateTime"`
	}
	type ProductionStage struct {
		IsReachEnd                 int    `json:"isReachEnd"`
		ProductionStageAwardStatus int    `json:"productionStageAwardStatus"`
		ProductionStageProgress    string `json:"productionStageProgress"`
	}
	type Speciality struct {
		FactoryFlowerQualification int `json:"factoryFlowerQualification"`
		FactoryFlowerStatus        int `json:"factoryFlowerStatus"`
		SkinQualification          int `json:"skinQualification"`
		SkinStatus                 int `json:"skinStatus"`
	}
	type User struct {
		CreateTime                int    `json:"createTime"`
		CurrentLevel              int    `json:"currentLevel"`
		DataMark                  int    `json:"dataMark"`
		DeviceID                  string `json:"deviceId"`
		Electric                  int    `json:"electric"`
		EncryptPin                string `json:"encryptPin"`
		HeadImage                 string `json:"headImage"`
		HongBaoValue              string `json:"hongBaoValue"`
		IsJXNewUser               int    `json:"isJXNewUser"`
		IsProductSpecialCommodity int    `json:"isProductSpecialCommodity"`
		MosaicPin                 string `json:"mosaicPin"`
		NewPlayerWelfareFlag      int    `json:"newPlayerWelfareFlag"`
		NextLevelPercent          int    `json:"nextLevelPercent"`
		Nickname                  string `json:"nickname"`
		NpcStep                   int    `json:"npcStep"`
		Pin                       string `json:"pin"`
		ShareQywx                 string `json:"shareQywx"`
		UpdateTime                int    `json:"updateTime"`
		UserIdentity              string `json:"userIdentity"`
		Xid                       string `json:"xid"`
		Zone                      string `json:"zone"`
	}
	type UserAttrExtInfo struct {
		Electric              int `json:"electric"`
		InvestElectricLimDays int `json:"investElectricLimDays"`
		LastProduceInvestTime int `json:"lastProduceInvestTime"`
		ProductLimFlag        int `json:"productLimFlag"`
		RewardType            int `json:"rewardType"`
		UserType              int `json:"userType"`
	}
	type Data struct {
		AssistCondition             AssistCondition             `json:"assistCondition"`
		AssistMaterialTuanCondition AssistMaterialTuanCondition `json:"assistMaterialTuanCondition"`
		DeviceList                  []DeviceList                `json:"deviceList"`
		FactoryList                 []FactoryList               `json:"factoryList"`
		NeedSelectPickSite          int                         `json:"needSelectPickSite"`
		NewFactoryFlower            NewFactoryFlower            `json:"newFactoryFlower"`
		PickSiteInfo                PickSiteInfo                `json:"pickSiteInfo"`
		ProductionList              []ProductionList            `json:"productionList"`
		ProductionStage             ProductionStage             `json:"productionStage"`
		Speciality                  Speciality                  `json:"speciality"`
		SystemVersion               string                      `json:"systemVersion"`
		User                        User                        `json:"user"`
		UserAttrExtInfo             UserAttrExtInfo             `json:"userAttrExtInfo"`
	}
	type AutoGenerated struct {
		Data    Data   `json:"data"`
		Msg     string `json:"msg"`
		NowTime int    `json:"nowTime"`
		Ret     int    `json:"ret"`
	}
	url := "https://m.jingxi.com/dreamfactory/userinfo/GetUserInfo?zone=dream_factory&pin=&sharePin=&shareType=&materialTuanPin=&materialTuanId=&needPickSiteInfo=1&source=&_time=1637631683565&_ts=1637631683565&timeStamp=&_stk=_time,_ts,materialTuanId,materialTuanPin,needPickSiteInfo,pin,sharePin,shareType,source,timeStamp,zone&_ste=1&_=1637631683575&sceneval=2&g_login_type=1&g_ty=ls"

	req := httplib.Get(url)
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Cookie", cookie)
	req.Header("User-Agent", "jdpingou;"+ua())
	req.Header("Accept-Language", "zh-cn")
	req.Header("Referer", "https://st.jingxi.com/pingou/dream_factory/index.html?ptag=7155.9.46")
	// req.Header("Accept-Encoding", "gzip, deflate, br")

	data, _ := req.Bytes()
	a := &AutoGenerated{}
	json.Unmarshal(data, a)

	desc := ""
	not := true
	if state != nil {
		not = false
	}
	if len(a.Data.ProductionList) > 0 && len(a.Data.FactoryList) > 0 {
		var production = a.Data.ProductionList[0]
		if production.InvestedElectric >= production.NeedElectric {
			if production.ExchangeStatus == 1 {
				desc = "可以兑换商品了"
			}
			if production.ExchangeStatus == 3 {
				not = false
				desc = "商品兑换已超时，请选择新商品进行制造"
			}
			// await exchangeProNotify()
		} else {
			not = false
			desc = fmt.Sprintf(`预计最快还需%d天生产完毕`, (production.NeedElectric-production.InvestedElectric)/(2*60*60*24))

		}
	} else {
		if len(a.Data.FactoryList) == 0 {
			not = false
			desc = "请手动开启活动"
		} else if len(a.Data.ProductionList) == 0 {
			not = false
			desc = "请手动选购商品进行生产"
		}
	}
	if desc == "" {
		not = false
	}
	desc += "🏭"
	if state != nil {
		state <- desc
	}
	if not {
		a叉哦叉哦(core.FetchCookieValue("pt_pin", cookie), "京喜工厂", desc)
	}
}

func jdsy(cookie string, desc chan string) {
	type AutoGenerated struct {
		Success bool        `json:"success"`
		Code    interface{} `json:"code"`
		BsCode  interface{} `json:"bsCode"`
		Message interface{} `json:"message"`
		Data    struct {
			List []struct {
				ActID     int         `json:"actId"`
				ReportID  interface{} `json:"reportId"`
				OrderID   interface{} `json:"orderId"`
				ApplyTime int64       `json:"applyTime"`
				Status    int         `json:"status"`
				TrialImg  string      `json:"trialImg"`
				TrialName string      `json:"trialName"`
				Text      struct {
					ID   int    `json:"id"`
					Text string `json:"text"`
				} `json:"text"`
				LeftTime      int `json:"leftTime"`
				TryButtonList []struct {
					ID     int         `json:"id"`
					Text   string      `json:"text"`
					Schema interface{} `json:"schema"`
				} `json:"tryButtonList"`
				Bottom           interface{}   `json:"bottom"`
				ActType          int           `json:"actType"`
				TaskType         int           `json:"taskType"`
				SkuID            string        `json:"skuId"`
				OrderState       interface{}   `json:"orderState"`
				Tag              []interface{} `json:"tag"`
				OrderAmount      interface{}   `json:"orderAmount"`
				EndTime          int64         `json:"endTime"`
				SupplierDelivery bool          `json:"supplierDelivery"`
			} `json:"list"`
			PageSize int   `json:"pageSize"`
			Page     int   `json:"page"`
			SysDate  int64 `json:"sysDate"`
		} `json:"data"`
	}
	rt := ""
	warn := make(chan string)
	go func() {
		req := httplib.Post("https://api.m.jd.com/client.action")
		req.Header("Host", "api.m.jd.com")
		req.Header("Content-Type", "application/x-www-form-urlencoded")
		req.Header("Origin", "https://prodev.m.jd.com")
		// req.Header("Accept-Encoding", "gzip, deflate, br")
		req.Header("Cookie", cookie)
		req.Header("Connection", "keep-alive")
		req.Header("Accept", "application/json, text/plain, */*")
		req.Header("User-Agent", ua())
		// req.Header("Referer", "https://prodev.m.jd.com/mall/active/2Y2YgUu1Xbbv8AfN7TAHhNqfQrAV/index.html?tttparams=eliIVi1eyJncHNfYXJlYSI6IjEyXzkzOV8yMzY4M181NjE4NCIsInByc3RhdGUiOiIwIiwidW5fYXJlYSI6IjEyXzkzOV8yMzY4M181NjE4NCIsIm1vZGVsIjoiaVBob25lMTAsMiIsImdMYXQiOiIzMy4yOTQ5MyIsImdMbmciOiIxMjAuMTQ4MjMyIiwibG5nIjoiMTIwLjE1MDMyMyIsImxhdCI6IjMzLjI5NTUzNi7J9&sid=a2e1c9b3b215a2337517cd01f2e04cbw&un_area=12_939_23683_56184")
		// req.Header("Content-Length", "332")
		// req.Header("Accept-Language", "zh-cn")
		req.Body(`appid=newtry&functionId=try_MyTrials&uuid=3345ad3d16ab2153c69f8ca91cd3e931b06a3bb8&clientVersion=10.2.7&client=wh5&osVersion=14.7.1&area=12_939_23683_56184&networkType=wifi&body=%7B%22geo%22%3A%7B%22lng%22%3A121.15326252577907%2C%22lat%22%3A34.295611038697575%7D%2C%22page%22%3A1%2C%22selected%22%3A2%2C%22previewTime%22%3A%22%22%7D`)
		//appid=newtry&functionId=try_MyTrials&uuid=3345ad3d16ab2153c69f8ca91cd3e931b06a3bb8&clientVersion=10.2.7&client=wh5&osVersion=14.7.1&area=12_939_23683_56184&networkType=wifi&body=%7B%22geo%22%3A%7B%22lng%22%3A121.15326252577907%2C%22lat%22%3A34.295611038697575%7D%2C%22page%22%3A1%2C%22selected%22%3A1%2C%22previewTime%22%3A%22%22%7D
		data, _ := req.Bytes()
		// fmt.Println(string(data))
		a := &AutoGenerated{}
		json.Unmarshal(data, a)
		// fmt.Println(a)
		for _, v := range a.Data.List {
			if len(v.TryButtonList) == 2 {
				if v.TryButtonList[0].ID <= 2 {
					warn <- "你有一个商品待领取，详情：" + v.TrialName + "👆"
					return
				}
			}
		}
		warn <- ""
	}()
	if rt == "" {
		req := httplib.Post("https://api.m.jd.com/client.action")
		req.Header("Host", "api.m.jd.com")
		req.Header("Content-Type", "application/x-www-form-urlencoded")
		req.Header("Origin", "https://prodev.m.jd.com")
		// req.Header("Accept-Encoding", "gzip, deflate, br")
		req.Header("Cookie", cookie)
		req.Header("Connection", "keep-alive")
		req.Header("Accept", "application/json, text/plain, */*")
		req.Header("User-Agent", ua())
		// req.Header("Referer", "https://prodev.m.jd.com/mall/active/2Y2YgUu1Xbbv8AfN7TAHhNqfQrAV/index.html?tttparams=eliIVi1eyJncHNfYXJlYSI6IjEyXzkzOV8yMzY4M181NjE4NCIsInByc3RhdGUiOiIwIiwidW5fYXJlYSI6IjEyXzkzOV8yMzY4M181NjE4NCIsIm1vZGVsIjoiaVBob25lMTAsMiIsImdMYXQiOiIzMy4yOTQ5MyIsImdMbmciOiIxMjAuMTQ4MjMyIiwibG5nIjoiMTIwLjE1MDMyMyIsImxhdCI6IjMzLjI5NTUzNi7J9&sid=a2e1c9b3b215a2337517cd01f2e04cbw&un_area=12_939_23683_56184")
		// req.Header("Content-Length", "332")
		// req.Header("Accept-Language", "zh-cn")
		req.Body(`appid=newtry&functionId=try_MyTrials&uuid=3345ad3d16ab2153c69f8ca91cd3e931b06a3bb8&clientVersion=10.2.7&client=wh5&osVersion=14.7.1&area=12_939_23683_56184&networkType=wifi&body=%7B%22geo%22%3A%7B%22lng%22%3A121.15326252577907%2C%22lat%22%3A34.295611038697575%7D%2C%22page%22%3A1%2C%22selected%22%3A1%2C%22previewTime%22%3A%22%22%7D`)

		data, _ := req.Bytes()
		// fmt.Println(string(data))
		a := &AutoGenerated{}
		json.Unmarshal(data, a)
		// fmt.Println(a)
		rt = fmt.Sprintf("%d件商品申请中", len(a.Data.List))
	}
	if xx := <-warn; xx != "" {
		rt = xx
	}
	desc <- rt
}

func cwwjf(cookie string, desc chan int) {
	req := httplib.Post("https://api.m.jd.com/api?appid=jdchoujiang_h5&functionId=giftGetBeanConfigs&body={%22reqSource%22:%22h5%22}")
	req.Header("Host", "api.m.jd.com")
	// req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Origin", "https://h5.m.jd.com")
	// req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Cookie", cookie)
	// req.Header("Connection", "keep-alive")
	req.Header("Accept", "application/json, text/plain, */*")
	req.Header("User-Agent", ua())
	// req.Header("Referer", "https://prodev.m.jd.com/mall/active/2Y2YgUu1Xbbv8AfN7TAHhNqfQrAV/index.html?tttparams=eliIVi1eyJncHNfYXJlYSI6IjEyXzkzOV8yMzY4M181NjE4NCIsInByc3RhdGUiOiIwIiwidW5fYXJlYSI6IjEyXzkzOV8yMzY4M181NjE4NCIsIm1vZGVsIjoiaVBob25lMTAsMiIsImdMYXQiOiIzMy4yOTQ5MyIsImdMbmciOiIxMjAuMTQ4MjMyIiwibG5nIjoiMTIwLjE1MDMyMyIsImxhdCI6IjMzLjI5NTUzNi7J9&sid=a2e1c9b3b215a2337517cd01f2e04cbw&un_area=12_939_23683_56184")
	// req.Header("Content-Length", "332")
	// req.Header("Accept-Language", "zh-cn")
	// req.Body(`appid=newtry&functionId=try_MyTrials&uuid=3345ad3d16ab2153c69f8ca91cd3e931b06a3bb8&clientVersion=10.2.7&client=wh5&osVersion=14.7.1&area=12_939_23683_56184&networkType=wifi&body=%7B%22geo%22%3A%7B%22lng%22%3A121.15326252577907%2C%22lat%22%3A34.295611038697575%7D%2C%22page%22%3A1%2C%22selected%22%3A1%2C%22previewTime%22%3A%22%22%7D`)
	data, _ := req.Bytes()
	coin, _ := jsonparser.GetInt(data, "data", "petCoin")
	desc <- int(coin)
}
