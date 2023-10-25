import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { CallbackComponent } from './callback/callback.component';
import { HomeComponent } from './home/home.component'; // if you have it

const routes: Routes = [
  { path: 'callback', component: CallbackComponent },
  { path: '', redirectTo: '/home', pathMatch: 'full' }, // assuming you want to redirect to '/home'
  { path: 'home', component: HomeComponent }, // your actual home route
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
