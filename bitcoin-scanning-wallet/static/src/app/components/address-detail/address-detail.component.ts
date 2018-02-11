import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
import { ApiService } from '../../services/api.service';

@Component({
  selector: 'app-home',
  styleUrls: ['./address-detail.component.css'],
  templateUrl: './address-detail.component.html',
})

export class AddressComponent {
  address: any;
  routeSubscription: Subscription;

  rows: Array< { txHash: string, blockHash: string, blockHeight: number, amount: number}> = [];
  columns = [
    { prop: 'txHash' },
    { prop: 'blockHash' },
    { prop: 'blockHeight' },
    { prop: 'amount' }
  ];

  constructor(private api: ApiService, private route: ActivatedRoute ) {
        this.routeSubscription = route.params.subscribe((params: any) => {
          this.api.GetAddress(params['id']).subscribe((data: any) => {
            console.log(data._body);
            let resp = JSON.parse(data._body);
            this.address = resp;
            this.changeRows(resp.txs);
          });
            });
  }


  public changeRows(txs: any) {
    this.rows = [];
    for (let i = 0; i <= txs.length - 1; i++) {
      console.log(txs[i].tx_hash);
      this.rows.push({'txHash': txs[i].tx_hash, 'blockHash': txs[i].block_hash, 'blockHeight': txs[i].block_height, 'amount': txs[i].bitcoin_amount});
    }
  }

}
