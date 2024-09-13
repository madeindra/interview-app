export namespace model {
	
	export class Chat {
	    text: string;
	    audio: string;
	
	    static createFrom(source: any = {}) {
	        return new Chat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.audio = source["audio"];
	    }
	}
	export class AnswerChatResponse {
	    language: string;
	    prompt?: Chat;
	    answer?: Chat;
	
	    static createFrom(source: any = {}) {
	        return new AnswerChatResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.prompt = this.convertValues(source["prompt"], Chat);
	        this.answer = this.convertValues(source["answer"], Chat);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class StartChatResponse {
	    id: string;
	    secret: string;
	    language: string;
	    text: string;
	    audio: string;
	
	    static createFrom(source: any = {}) {
	        return new StartChatResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.secret = source["secret"];
	        this.language = source["language"];
	        this.text = source["text"];
	        this.audio = source["audio"];
	    }
	}
	export class StatusResponse {
	    server: boolean;
	    key: boolean;
	    api?: boolean;
	    apiStatus: string;
	
	    static createFrom(source: any = {}) {
	        return new StatusResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.key = source["key"];
	        this.api = source["api"];
	        this.apiStatus = source["apiStatus"];
	    }
	}

}

