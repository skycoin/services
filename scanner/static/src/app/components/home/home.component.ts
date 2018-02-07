import { Component } from '@angular/core';
import { ApiService } from '../../services/api.service';

@Component({
  selector: 'app-home',
  styleUrls: ['./home.component.css'],
  templateUrl: './home.component.html',
})

export class HomeComponent {
  addrs: any;
  new_addrs: '';

  rows: Array< { currency: string, address: string, transactions: string}> = [];
  columns = [
     { prop: 'currency' },
     { prop: 'address' },
     { prop: 'transactions' },
    ];

  constructor (private apiService: ApiService){
    this.getAddrs();
  }

  getAddrs() {
    this.apiService.GetWallet().subscribe((data: any) => {
      this.addrs = JSON.parse(data._body);
      console.log(this.addrs);
      this.changeRows(this.addrs);
    });
  }

  sendAddrs() {
    let split_addrs = this.new_addrs.split('\n');
    console.log(split_addrs)
    this.apiService.AddAddresses(split_addrs).subscribe((data) => {
      console.log(data);
      this.getAddrs();
    });

  }

  scanMin() {
    this.apiService.ScanMin().subscribe((data: any) => {
      this.addrs = JSON.parse(data._body);
      console.log(this.addrs);
      this.changeRows(this.addrs);
    });
  }

  scanMax() {
    this.apiService.ScanMax().subscribe((data: any) => {
      this.addrs = JSON.parse(data._body);
      console.log(this.addrs);
      this.changeRows(this.addrs);
    });
  }

  scanFar() {
    this.apiService.ScanFar().subscribe((data: any) => {
      this.addrs = JSON.parse(data._body);
      console.log(this.addrs);
      this.changeRows(this.addrs);
    });
  }

  scanShort() {
    this.apiService.ScanShort().subscribe((data: any) => {
      this.addrs = JSON.parse(data._body);
      console.log(this.addrs);
      this.changeRows(this.addrs);
    });
  }

  public changeRows(addresses: any) {
    this.rows = [];
      for (let i = 0; i <= addresses.length - 1; i++) {
          this.rows.push({'currency': 'Bitcoin', 'address': addresses[i].address, 'transactions': addresses[i].txs.length});
       }
  }

}
