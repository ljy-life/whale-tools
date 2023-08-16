/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ljy-life/whale-tools.git/config"
	"os"
	"os/signal"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping [-c count] [-i interval] [-t timeout] [--privileged] host",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logger.Error("ping need host")
			return
		}
		if len(args) > 1 {
			logger.Error("ping only support one host")
			return
		}
		host := args[0]
		Ping(host)
	},
}

var logger = config.Log

// PingConfig ping 配置
type PingConfig struct {
	Count      int
	Size       int
	TTl        int
	Interval   time.Duration
	Timeout    time.Duration
	Privileged bool
}

var pingConfig PingConfig

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().IntVarP(&pingConfig.Count, "count", "c", -1, " 停止 ping 的次数")
	pingCmd.Flags().IntVarP(&pingConfig.Size, "size", "s", 24, "发送 ICMP 消息时的负载大小")
	pingCmd.Flags().IntVarP(&pingConfig.TTl, "ttl", "", 64, "设置 IP 数据包的 TTL 字段，即生存时间")
	pingCmd.Flags().DurationVarP(&pingConfig.Interval, "interval", "i", time.Second, "等待两次 ping 之间的时间间隔")
	pingCmd.Flags().DurationVarP(&pingConfig.Timeout, "timeout", "t", time.Second*100000, "停止 ping 的时间")
	pingCmd.Flags().BoolVarP(&pingConfig.Privileged, "privileged", "", false, "发送特权原始 ICMP ping")
	pingCmd.Flags().BoolP("help", "h", false, "查看帮助")
}

func Ping(host string) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		logger.Errorf("ping %s error: %s", host, err)
		return
	}
	//  监听 Ctrl+C 退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
		}
	}()
	pinger.OnRecv = func(pkt *probing.Packet) {
		logger.Infof("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}
	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		logger.Infof("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}
	pinger.OnFinish = func(stats *probing.Statistics) {
		logger.Infof("\n--- %s ping statistics ---\n", stats.Addr)
		logger.Infof("%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
		logger.Infof("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}
	pinger.Count = pingConfig.Count
	pinger.Size = pingConfig.Size
	pinger.Interval = pingConfig.Interval
	pinger.Timeout = pingConfig.Timeout
	pinger.TTL = pingConfig.TTl
	pinger.SetPrivileged(pingConfig.Privileged)
	logger.Infof("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		logger.Errorf("Failed to ping target host:", err)
	}
}
