import { Routes } from '@angular/router';

import { HomeComponent } from './components/home/home.component';
import { AddressComponent } from './components/address-detail/address-detail.component';



export const AppRoutes: Routes = [
  {
    path: '',
    component: HomeComponent
  },
  {
     path: 'address/:id',
     component: AddressComponent,
     data: { title: '' }
 },
]
