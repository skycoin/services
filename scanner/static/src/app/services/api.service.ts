import { Injectable } from '@angular/core';
import { Headers, Http, RequestOptions } from '@angular/http';
import 'rxjs/add/operator/map';

@Injectable()
export class ApiService {
  private api: string = 'http://localhost:7755'


  constructor (private http: Http){
  }

  public GetWallet() {
    return this.http.get(this.api + '/getaddrs');
  }

  public GetAddress(address: string) {
    return this.http.get(this.api + '/getaddr?address=' + address);
  }

  public ScanMin() {
    return this.http.get(this.api + '/scanmin');
  }

  public ScanMax() {
    return this.http.get(this.api + '/scanmax');
  }

  public ScanFar() {
    return this.http.get(this.api + '/scanfar');
  }

  public ScanShort() {
    return this.http.get(this.api + '/scanshort');
  }

  public AddAddresses(Addrs: string[]) {
    return this.http.post(this.api + '/newaddrs', JSON.stringify({addrs: Addrs}));
  }


}
