export namespace main {
	
	export class StreamStatus {
	    isStreaming: boolean;
	    videoPath: string;
	    rtmpUrl: string;
	    error: string;
	    elapsedSeconds: number;
	    durationSeconds: number;
	    remainingSeconds: number;
	    quality: string;
	    connectionHealth: string;
	    retryCount: number;
	    maxRetries: number;
	
	    static createFrom(source: any = {}) {
	        return new StreamStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isStreaming = source["isStreaming"];
	        this.videoPath = source["videoPath"];
	        this.rtmpUrl = source["rtmpUrl"];
	        this.error = source["error"];
	        this.elapsedSeconds = source["elapsedSeconds"];
	        this.durationSeconds = source["durationSeconds"];
	        this.remainingSeconds = source["remainingSeconds"];
	        this.quality = source["quality"];
	        this.connectionHealth = source["connectionHealth"];
	        this.retryCount = source["retryCount"];
	        this.maxRetries = source["maxRetries"];
	    }
	}

}

